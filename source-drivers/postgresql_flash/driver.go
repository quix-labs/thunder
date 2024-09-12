package postgresql_flash

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/quix-labs/thunder"
	"strconv"
	"time"
)

func init() {
	thunder.RegisterSourceDriver(&Driver{})
}

type DriverConfig struct {
	Host     string `default:"localhost"`
	Port     uint16 `type:"number" default:"5432"`
	User     string `required:"true"`
	Password string `required:"true" type:"password"`
	Database string `required:"true"`
	Schema   string `default:"public"`
}

type Driver struct {
	config *DriverConfig
	conn   *pgx.Conn
}

//go:embed icon.svg
var logo string

func (d *Driver) DriverInfo() thunder.SourceDriverInfo {
	return thunder.SourceDriverInfo{
		ID: "postgresql_flash",
		New: func(config any) (thunder.SourceDriver, error) {

			cfg, ok := config.(*DriverConfig)
			if !ok {
				return nil, errors.New("invalid config type")
			}

			pgConn, err := d.newConn(cfg)
			if err != nil {
				return nil, err
			}

			return &Driver{
				config: cfg,
				conn:   pgConn,
			}, nil
		},

		Name:   "PostgreSQL (Flash)",
		Image:  logo,
		Config: DriverConfig{},
		Notes:  []string{"Please be sure wal_replication is set to logical on your PostgreSQL config"},
	}
}

func (d *Driver) TestConfig() (string, error) {
	stats, err := d.Stats()
	if err != nil {
		return "cannot get field table statistics", err
	}
	return fmt.Sprintf("Successfully connected, discover %d tables", len(*stats)), nil
}

func (d *Driver) Stats() (*thunder.SourceDriverStats, error) {
	query := StatsQuery(d.config.Schema)
	stats := thunder.SourceDriverStats{}
	type RowResult struct {
		Name        string   `json:"name"`
		Columns     []string `json:"columns"`
		PrimaryKeys []string `json:"primary_keys"`
	}

	results, err := GetResultsSync[RowResult](d.conn, query, time.Second*10, false)
	if err != nil {
		return nil, err
	}
	for _, result := range results {
		stats[result.Name] = thunder.SourceDriverStatsTable{
			Columns:     result.Columns,
			PrimaryKeys: result.PrimaryKeys,
		}
	}
	return &stats, nil
}

// TODO GET CHAN AS ARGS AND RETURN ERROR
func (d *Driver) GetDocumentsForProcessor(processor *thunder.Processor, limit uint64) (<-chan *thunder.Document, <-chan error) {
	query, err := GetSqlForProcessor(processor)
	if err != nil {
		panic(err)
	}
	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT %s", query, strconv.FormatUint(limit, 10))
	}

	return GetResultsChan[thunder.Document](d.conn, query, false)
}

func (d *Driver) Start() error {
	//TODO implement me
	panic("implement me")
}

func (d *Driver) Stop() error {
	if d.conn != nil {
		return d.conn.Close(context.Background())
	}
	return nil
}

func (d *Driver) Shutdown() error {
	if d.conn != nil {
		return d.conn.Close(context.Background())
	}
	return nil
}

func (d *Driver) newConn(cfg *DriverConfig) (*pgx.Conn, error) {
	pgConnConfig, err := pgx.ParseConfig("")
	if err != nil {
		return nil, err
	}
	pgConnConfig.Host = cfg.Host
	pgConnConfig.User = cfg.User
	pgConnConfig.Port = cfg.Port
	pgConnConfig.Password = cfg.Password
	pgConnConfig.Database = cfg.Database

	pgConn, err := pgx.ConnectConfig(context.Background(), pgConnConfig)
	if err != nil {
		return nil, err
	}
	return pgConn, nil
}

var (
	_ thunder.SourceDriver = (*Driver)(nil) // Interface implementation
)
