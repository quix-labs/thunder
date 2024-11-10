package elastic

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/quix-labs/thunder"
	"strings"
	"sync/atomic"
	"time"
)

const DriverID = "thunder.elastic"

func init() {
	_ = thunder.TargetDrivers.Register(DriverID, &Driver{})
}

type Driver struct {
	config *DriverConfig
	client *elasticsearch.TypedClient
}

func (d *Driver) New(config any) (thunder.TargetDriver, error) {
	cfg, ok := config.(*DriverConfig)
	if !ok {
		return nil, errors.New("invalid config type")
	}

	client, err := NewConn(cfg)
	if err != nil {
		return nil, err
	}

	return &Driver{config: cfg, client: client}, nil
}

func (d *Driver) ID() string {
	return DriverID
}

//go:embed icon.svg
var logo string

func (d *Driver) Config() thunder.TargetDriverConfig {
	return thunder.TargetDriverConfig{
		Name:   "Elasticsearch",
		Config: DriverConfig{},
		Image:  logo,
		Notes:  nil,
	}
}

func (d *Driver) TestConfig() (string, error) {
	res, err := d.client.Info().Do(context.Background())
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`successfully connected, cluster: "%s"`, res.ClusterName), nil
}

func (d *Driver) HandleEvents(processor *thunder.Processor, eventsChan <-chan thunder.TargetEvent) error {
	bulkIndexer := NewBulkIndexer(d.client, int64(d.config.BatchMaxBytesSize*1024), int64(d.config.ParallelBatch))
	defer bulkIndexer.Close()
	bulkIndexer.AddSendTimeout(time.Second * time.Duration(d.config.ReactivityInterval))

	index := d.config.Prefix + processor.Index
	var dataReceived atomic.Bool

	for event := range eventsChan {
		dataReceived.Store(true)
		if err := d.processEvent(event, bulkIndexer, index); err != nil {
			return err
		}
	}

	if dataReceived.Load() {
		bulkIndexer.Close()
		_, err := d.client.Indices.Flush().Index(index).Do(context.Background())
		if err != nil {
			return fmt.Errorf("failed to flush index: %w", err)
		}
	}
	return nil
}

func (d *Driver) processEvent(event thunder.TargetEvent, bulkIndexer *BulkIndexer, index string) error {
	switch typedEvent := event.(type) {
	case *thunder.TargetInsertEvent:
		return bulkIndexer.Add(
			[]byte(`{"index":{"_index":"`+index+`","_id":"`+SanitizeJsonString(typedEvent.Pkey)+`"}}`),
			typedEvent.Json,
		)
	case *thunder.TargetPatchEvent:
		// Handle base table update
		if typedEvent.Relation == nil {
			return bulkIndexer.Add(
				[]byte(`{"update":{"_index":"`+index+`","_id":"`+SanitizeJsonString(typedEvent.Pkey)+`"}}`),
				bytes.Join([][]byte{[]byte(`{"doc":`), typedEvent.JsonPatch, []byte("}")}, []byte("")),
			)
		}

		// Handle relation patch
		path := typedEvent.Relation.Path()
		var reqBody string
		if typedEvent.Relation.Many {
			reqBody = fmt.Sprintf(`{
						"script": {
							"source": "for (int i = 0; i < ctx._source.%s.size(); i++) { if (ctx._source.%s[i]._pkey == '%s') { for (entry in params.patch.entrySet()) { ctx._source.%s[i][entry.getKey()] = entry.getValue(); } } }",
							"params": {"patch": %s}
						},
						"query": {"match": {"%s._pkey": "%s"}}
					}`, path, path, SanitizeJsonString(typedEvent.Pkey), path, typedEvent.JsonPatch, path, SanitizeJsonString(typedEvent.Pkey))
		} else {
			reqBody = fmt.Sprintf(`{
						"script": {
							"source": "for (entry in params.patch.entrySet()) { ctx._source.%s[entry.getKey()] = entry.getValue(); }",
							"params": {"patch": %s}
						},
						"query": {"match": {"%s._pkey": "%s"}}
					}`, path, typedEvent.JsonPatch, path, SanitizeJsonString(typedEvent.Pkey))
		}

		var refresh = true // IMPORTANT TO AVOID BURST 409
		req := esapi.UpdateByQueryRequest{
			Index:   []string{index},
			Refresh: &refresh,
			Body:    strings.NewReader(reqBody),
		}

		res, err := req.Do(context.Background(), d.client)
		if err != nil {
			return err
		}

		if res.IsError() {
			err := errors.New(res.String())
			closeErr := res.Body.Close()
			return errors.Join(err, closeErr)
		}

		return res.Body.Close()

	case *thunder.TargetDeleteEvent:
		// Handle base table delete
		if typedEvent.Relation == nil {
			return bulkIndexer.Add(
				[]byte(`{"delete":{"_index":"` + index + `","_id":"` + SanitizeJsonString(typedEvent.Pkey) + `"}}`),
			)
		}

		// Handle relation delete
		path := typedEvent.Relation.Path()
		var reqBody string
		if typedEvent.Relation.Many == true {
			reqBody = fmt.Sprintf(`{
						"script": "ctx._source.%s.removeIf(item -> item._pkey =='%s')",
						"query": {"match": {"%s._pkey":"%s"}}
					}`, path, SanitizeJsonString(typedEvent.Pkey), path, SanitizeJsonString(typedEvent.Pkey))

		} else {
			reqBody = fmt.Sprintf(`{
						"script": "ctx._source.%s=null",
						"query": {"match": {"%s._pkey":"%s"}}
					}`, path, path, SanitizeJsonString(typedEvent.Pkey))
		}

		var refresh = true // IMPORTANT TO AVOID BURST 409
		req := esapi.UpdateByQueryRequest{Index: []string{index},
			Refresh: &refresh,
			Body:    strings.NewReader(reqBody),
		}

		res, err := req.Do(context.Background(), d.client)
		if err != nil {
			return err
		}

		if res.IsError() {
			err := errors.New(res.String())
			closeErr := res.Body.Close()
			return errors.Join(err, closeErr)
		}

		return res.Body.Close()

	case *thunder.TargetTruncateEvent:
		// Handle base table truncate
		if typedEvent.Relation == nil {
			_, err := d.client.DeleteByQuery(index).
				Query(&types.Query{MatchAll: &types.MatchAllQuery{}}).
				Refresh(true). // IMPORTANT TO AVOID BURST 409
				Do(context.Background())
			return err
		}

		// Handle relation truncate
		path := typedEvent.Relation.Path()

		scriptNewValue := "null"
		if typedEvent.Relation.Many {
			scriptNewValue = "[]"
		}

		reqBody := fmt.Sprintf(`{
					"script": {"source": "if (ctx._source.%s != null) { ctx._source.%s = %s; }"},
					"query": {"exists": {"field": "%s._pkey"}}
				}`, path, path, scriptNewValue, path)

		var refresh = true // IMPORTANT TO AVOID BURST 409
		req := esapi.UpdateByQueryRequest{
			Index:   []string{index},
			Refresh: &refresh,
			Body:    strings.NewReader(reqBody),
		}

		res, err := req.Do(context.Background(), d.client)
		if err != nil {
			return err
		}

		if res.IsError() {
			err := errors.New(res.String())
			closeErr := res.Body.Close()
			return errors.Join(err, closeErr)
		}

		return res.Body.Close()

	default:
		return fmt.Errorf("unsupported event type received: %T", typedEvent)
	}
}

func (d *Driver) Shutdown() error {
	return nil
}

func NewConn(cfg *DriverConfig) (*elasticsearch.TypedClient, error) {
	esConfig := elasticsearch.Config{}
	esConfig.Addresses = []string{cfg.Endpoint}
	esConfig.Username = cfg.Username
	esConfig.Password = cfg.Password
	return elasticsearch.NewTypedClient(esConfig)
}

var (
	_ thunder.TargetDriver = (*Driver)(nil)
)
