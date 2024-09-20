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
				if typedEvent.Path == "" {
					bulkIndexer.Add(
						[]byte(`{"update":{"_index":"`+index+`","_id":"`+SanitizeJsonString(typedEvent.Pkey)+`"}}`),
						bytes.Join([][]byte{[]byte(`{"doc":`), typedEvent.JsonPatch, []byte("}")}, []byte("")),
					)
					continue
				}

				// Handle relation patch
				reqBody := fmt.Sprintf(`
    {
        "script": {
            "source": "for (entry in params.patch.entrySet()) { ctx._source.%s[entry.getKey()] = entry.getValue(); }",
            "params": {
                "patch": %s
            }
        },
        "query": {
            "match": {
                "%s._pkey": "%s" 
            }
        }
    }`, typedEvent.Path, typedEvent.JsonPatch, typedEvent.Path, SanitizeJsonString(typedEvent.Pkey))
				//TODO MANY TO MANY ARRAY
				//"source": "for (int i = 0; i < ctx._source.%s.length; i++) { if (ctx._source.%s[i].author == params.oldAuthor) { ctx._source.%s[i].author = params.newAuthor; } }",

				req := esapi.UpdateByQueryRequest{
					Index: []string{index},
					Body:  strings.NewReader(reqBody),
				}

				res, err := req.Do(ctx, d.client)
				if err != nil {
					return err
				}
				if err = res.Body.Close(); err != nil {
					return err
				}

			case *thunder.TargetDeleteEvent:
				// Handle base table delete
				if typedEvent.Path == "" {
					bulkIndexer.Add(
						[]byte(`{"delete":{"_index":"` + index + `","_id":"` + SanitizeJsonString(typedEvent.Pkey) + `"}}`),
					)
					continue
				}

				// Handle relation delete
				reqBody := fmt.Sprintf(`{
					"script": {"source": "if (ctx._source.%s != null) { ctx._source.%s = null; }"},
					"query": {"match": {"%s._pkey":"%s"}}
				}`, typedEvent.Path, typedEvent.Path, typedEvent.Path, SanitizeJsonString(typedEvent.Pkey))

				//TODO MANY TO MANY ARRAY

				var refresh = true // IMPORTANT TO AVOID BURST 409
				req := esapi.UpdateByQueryRequest{Index: []string{index},
					Refresh: &refresh,
					Body:    strings.NewReader(reqBody),
				}

				res, err := req.Do(ctx, d.client)
				if err != nil {
					return err
				}

				if err = res.Body.Close(); err != nil {
					return err
				}

			case *thunder.TargetTruncateEvent:
				// Handle base table truncate
				if typedEvent.Path == "" {
					res, err := d.client.DeleteByQuery(index).
						Query(&types.Query{MatchAll: &types.MatchAllQuery{}}).
						Refresh(true). // IMPORTANT TO AVOID BURST 409
						Do(context.Background())
					if err != nil {
						fmt.Println(err)
						return err
					}
					fmt.Println(res.Deleted)
					continue
				}

				// Handle relation truncate
				reqBody := fmt.Sprintf(`{
					"script": {"source": "if (ctx._source.%s != null) { ctx._source.%s = null; }"},
					"query": {"exists": {"field": "%s._pkey"}}
				}`, typedEvent.Path, typedEvent.Path, typedEvent.Path)

				var refresh = true // IMPORTANT TO AVOID BURST 409
				req := esapi.UpdateByQueryRequest{Index: []string{index},
					Refresh: &refresh,
					Body:    strings.NewReader(reqBody),
				}

				res, err := req.Do(ctx, d.client)
				if err != nil {
					return err
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
