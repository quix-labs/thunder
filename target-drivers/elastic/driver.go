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
	"time"
)

func init() {
	thunder.RegisterTargetDriver(&Driver{})
}

type Driver struct {
	config *DriverConfig
	client *elasticsearch.TypedClient
}

//go:embed icon.svg
var logo string

func (d *Driver) DriverInfo() thunder.TargetDriverInfo {
	return thunder.TargetDriverInfo{
		ID: "elastic",
		New: func(config any) (thunder.TargetDriver, error) {
			cfg, ok := config.(*DriverConfig)
			if !ok {
				return nil, errors.New("invalid config type")
			}

			client, err := NewConn(cfg)
			if err != nil {
				return nil, err
			}

			return &Driver{config: cfg, client: client}, nil
		},
		Name:   "Elasticsearch",
		Config: DriverConfig{},
		Image:  logo,
	}
}

func (d *Driver) TestConfig() (string, error) {
	res, err := d.client.Info().Do(context.Background())
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`successfully connected, cluster: "%s"`, res.ClusterName), nil
}

func (d *Driver) HandleEvents(processor *thunder.Processor, eventsChan <-chan thunder.TargetEvent, ctx context.Context) error {
	bulkIndexer := NewBulkIndexer(d.client, d.config.BatchSize)
	defer bulkIndexer.Close()
	bulkIndexer.AddSendTimeout(time.Second * time.Duration(d.config.ReactivityInterval))

	// Defer flush index
	defer func() {
		_, _ = d.client.Indices.Flush().Index(processor.Index).Do(context.Background())
	}()

	index := d.config.Prefix + processor.Index
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case event := <-eventsChan:
			switch typedEvent := event.(type) {
			case *thunder.TargetInsertEvent:
				bulkIndexer.Add(
					[]byte(`{"index":{"_index":"`+index+`","_id":"`+SanitizeJsonString(typedEvent.Pkey)+`"}}`),
					typedEvent.Json,
				)
				break
			case *thunder.TargetPatchEvent:
				// Handle base table update
				if typedEvent.Relation == nil {
					bulkIndexer.Add(
						[]byte(`{"update":{"_index":"`+index+`","_id":"`+SanitizeJsonString(typedEvent.Pkey)+`"}}`),
						bytes.Join([][]byte{[]byte(`{"doc":`), typedEvent.JsonPatch, []byte("}")}, []byte("")),
					)
					continue
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

				res, err := req.Do(ctx, d.client)
				if err != nil {
					return err
				}

				if res.IsError() {
					err := errors.New(res.String())
					_ = res.Body.Close()
					fmt.Println(err)
					return err
				}

				if err = res.Body.Close(); err != nil {
					return err
				}

			case *thunder.TargetDeleteEvent:
				// Handle base table delete
				if typedEvent.Relation == nil {
					bulkIndexer.Add(
						[]byte(`{"delete":{"_index":"` + index + `","_id":"` + SanitizeJsonString(typedEvent.Pkey) + `"}}`),
					)
					continue
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

				res, err := req.Do(ctx, d.client)
				if err != nil {
					return err
				}

				if res.IsError() {
					_ = res.Body.Close()
					return fmt.Errorf(res.String())
				}

				if err = res.Body.Close(); err != nil {
					return err
				}

			case *thunder.TargetTruncateEvent:
				// Handle base table truncate
				if typedEvent.Relation == nil {
					_, err := d.client.DeleteByQuery(index).
						Query(&types.Query{MatchAll: &types.MatchAllQuery{}}).
						Refresh(true). // IMPORTANT TO AVOID BURST 409
						Do(context.Background())
					if err != nil {
						fmt.Println(err)
						return err
					}
					continue
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

				res, err := req.Do(ctx, d.client)
				if err != nil {
					return err
				}

				if res.IsError() {
					_ = res.Body.Close()
					return fmt.Errorf(res.String())
				}

				if err = res.Body.Close(); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unsupported event type received: %T", typedEvent)
			}
		}
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
