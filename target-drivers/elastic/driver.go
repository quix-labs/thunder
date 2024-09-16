package elastic

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/quix-labs/thunder"
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

func (d *Driver) HandleEvents(processor *thunder.Processor, eventsChan <-chan thunder.Event, ctx context.Context) error {
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
			case *thunder.InsertEvent:
				bulkIndexer.Add(
					[]byte(`{"index":{"_index":"`+index+`","_id":"`+GetPrimaryKeysAsString(typedEvent.PrimaryKeys)+`"}}`),
					typedEvent.Json,
				)
				break
			case *thunder.PatchEvent:
				//TODO PATH
				bulkIndexer.Add(
					[]byte(`{"update":{"_index":"`+index+`","_id":"`+GetPrimaryKeysAsString(typedEvent.PrimaryKeys)+`"}}`),
					bytes.Join([][]byte{[]byte(`{"doc":`), typedEvent.JsonPatch, []byte("}")}, []byte("")),
				)
			case *thunder.DeleteEvent:
				bulkIndexer.Add(
					[]byte(`{"delete":{"_index":"` + index + `","_id":"` + GetPrimaryKeysAsString(typedEvent.PrimaryKeys) + `"}}`),
				)
			case *thunder.TruncateEvent:
				// not supported by bulk
				_, err := d.client.DeleteByQuery(index).
					Query(&types.Query{MatchAll: &types.MatchAllQuery{}}).
					Do(context.Background())
				if err != nil {
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
