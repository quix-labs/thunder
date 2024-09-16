package elastic

import (
	"bytes"
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"sync"
	"sync/atomic"
	"time"
)

type BulkIndexer struct {
	mu      sync.Mutex
	_queue  [][]byte
	_client *elasticsearch.TypedClient

	_lastSent         time.Time
	pendingOperations atomic.Int64
	BatchMaxSize      int
}

func NewBulkIndexer(client *elasticsearch.TypedClient, batchMaxSize int) *BulkIndexer {
	return &BulkIndexer{
		_queue:       make([][]byte, 0, batchMaxSize*2),
		_client:      client,
		BatchMaxSize: batchMaxSize,
	}
}

func (bi *BulkIndexer) Add(lines ...[]byte) {
	bi.mu.Lock()
	bi._queue = append(bi._queue, lines...)
	bi.mu.Unlock()

	bi.pendingOperations.Add(1)

	if bi.pendingOperations.Load() >= int64(bi.BatchMaxSize) {
		bi.flush()
	}
}

func (bi *BulkIndexer) Close() {
	if bi.pendingOperations.Load() > 0 {
		bi.flush()
	}
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

	fullBody := append(bytes.Join(bi._queue, []byte("\n")), "\n"...)
	_, err := bi._client.Bulk().Raw(bytes.NewReader(fullBody)).Do(context.Background())
	if err != nil {
		panic(err)
		return
	}

	bi._lastSent = time.Now()
	bi.pendingOperations.Store(int64(0))
	bi._queue = bi._queue[:0]
}
