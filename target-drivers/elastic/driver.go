package elastic

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/quix-labs/thunder"
)

func init() {
	thunder.RegisterTargetDriver(&Driver{})
}

type DriverConfig struct {
	Endpoint  string `default:"http://localhost:9200"`
	Username  string
	Password  string `type:"password"`
	BatchSize int    `type:"number" label:"Batch size" default:"100" min:"1" help:"Use 1 to disable batching (not recommended)"`
	Prefix    string
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

func NewConn(cfg *DriverConfig) (*elasticsearch.TypedClient, error) {
	esConfig := elasticsearch.Config{}
	esConfig.Addresses = []string{cfg.Endpoint}
	esConfig.Username = cfg.Username
	esConfig.Password = cfg.Password
	return elasticsearch.NewTypedClient(esConfig)
}
