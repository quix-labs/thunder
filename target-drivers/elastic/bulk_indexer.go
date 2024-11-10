package elastic

import (
	"bytes"
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"golang.org/x/sync/errgroup"
	"sync"
	"sync/atomic"
	"time"
)

type BulkIndexer struct {
	mu sync.Mutex

	_queue [][][]byte

	_client *elasticsearch.TypedClient

	_lastSent time.Time

	pendingOperations atomic.Int64
	BatchMaxBytesSize int64
	ParallelBatch     int64

	_incrementFn func(lines *[][]byte) int64
	_closeOnce   sync.Once
}

func NewBulkIndexer(client *elasticsearch.TypedClient, batchMaxBytesSize int64, parallelBatch int64) *BulkIndexer {
	return &BulkIndexer{
		_queue: make([][][]byte, 0), // Dynamic allocation depends on item size

		_client:           client,
		BatchMaxBytesSize: batchMaxBytesSize,
		ParallelBatch:     parallelBatch,
		_incrementFn: func(lines *[][]byte) int64 {
			var size int64
			for _, line := range *lines {
				size += int64(len(line))
			}
			return size
		},
	}
}

func (bi *BulkIndexer) Add(lines ...[]byte) {
	bi.mu.Lock()
	bi._queue = append(bi._queue, lines)
	bi.mu.Unlock()

	bi.pendingOperations.Add(bi._incrementFn(&lines))

	if bi.pendingOperations.Load() >= (bi.BatchMaxBytesSize * bi.ParallelBatch) {
		bi.flush()
	}
}

func (bi *BulkIndexer) Close() {
	bi._closeOnce.Do(func() {
		if bi.pendingOperations.Load() > 0 {
			bi.flush()
		}
	})
}

func (bi *BulkIndexer) AddSendTimeout(duration time.Duration) {
	go func() {
		ticker := time.NewTicker(duration / 2)
		defer ticker.Stop()

		for {
			<-ticker.C

			bi.mu.Lock()
			timeSinceLastSend := time.Since(bi._lastSent)
			pendingOps := bi.pendingOperations.Load()
			bi.mu.Unlock()

			if pendingOps > 0 && timeSinceLastSend >= duration {
				bi.flush()
			}
		}
	}()
}

func (bi *BulkIndexer) flush() {
	bi.mu.Lock()
	defer bi.mu.Unlock()

	var eg = new(errgroup.Group)

	buffer := make([]byte, 0, int(bi.BatchMaxBytesSize))

	for _, action := range bi._queue {
		actionBody := append(bytes.Join(action, []byte("\n")), '\n')

		if len(buffer)+len(actionBody) > int(bi.BatchMaxBytesSize) {
			// Create local copy of buffer
			bodyToSend := make([]byte, len(buffer))
			copy(bodyToSend, buffer)

			eg.Go(func() error {
				_, err := bi._client.Bulk().Raw(bytes.NewReader(bodyToSend)).Do(context.Background())
				bodyToSend = nil // Force Golang GC to clean memory
				return err
			})
			buffer = buffer[:0]
		}

		buffer = append(buffer, actionBody...)
	}

	if len(buffer) > 0 {
		eg.Go(func() error {
			_, err := bi._client.Bulk().Raw(bytes.NewReader(buffer)).Do(context.Background())
			return err
		})
	}

	err := eg.Wait()
	if err != nil {
		panic(err)
		return
	}

	bi._lastSent = time.Now()
	bi.pendingOperations.Store(int64(0))

	bi._queue = bi._queue[:0]
}
