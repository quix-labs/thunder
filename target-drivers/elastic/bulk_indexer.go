package elastic

import (
	"bytes"
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"sync"
)

type BulkIndexer struct {
	mu      sync.Mutex
	_queue  [][]byte
	_client *elasticsearch.TypedClient

	BatchMaxSize int
}

func NewBulkIndexer(client *elasticsearch.TypedClient, batchMaxSize int) *BulkIndexer {
	return &BulkIndexer{
		_queue:       make([][]byte, 0, batchMaxSize*2),
		_client:      client,
		BatchMaxSize: batchMaxSize,
	}
}

func (bi *BulkIndexer) Add(operation []byte, data []byte) {
	bi.mu.Lock()
	bi._queue = append(bi._queue, operation, data)
	bi.mu.Unlock()

	if len(bi._queue) >= bi.BatchMaxSize*2 {
		bi.flush()
	}
}

func (bi *BulkIndexer) Close() {
	if len(bi._queue) > 0 {
		bi.flush()
	}
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

	bi._queue = bi._queue[:0]
}
