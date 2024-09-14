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

type Driver struct {
	config *DriverConfig
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
			return &Driver{config: cfg}, nil
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
	conn, err := d.newConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	query := StatsQuery(d.config.Schema)
	stats := thunder.SourceDriverStats{}
	type RowResult struct {
		Name        string   `json:"name"`
		Columns     []string `json:"columns"`
		PrimaryKeys []string `json:"primary_keys"`
	}

	results, err := GetResultsSync[RowResult](conn, query, time.Second*10, false)
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

func (d *Driver) GetDocumentsForProcessor(processor *thunder.Processor, docChan chan<- *thunder.Document, errChan chan error, limit uint64) {
	conn, err := d.newConn()
	if err != nil {
		errChan <- err
		return
	}
	defer conn.Close(context.Background())

	query, err := GetSqlForProcessor(processor)
	if err != nil {
		errChan <- err
		return
	}
	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT %s", query, strconv.FormatUint(limit, 10))
	}

	GetResultsInChan[thunder.Document](conn, query, false, docChan, errChan)
}

func (d *Driver) Start() error {
	//TODO implement me
	panic("implement me")
}

func (d *Driver) Stop() error {

	return nil
}

func (d *Driver) Shutdown() error {

	return nil
}

func (d *Driver) newConn() (*pgx.Conn, error) {
	pgConnConfig, err := pgx.ParseConfig("postgres://u:s@l:5432/d?sslmode=disable")
	if err != nil {
		return nil, err
	}
	pgConnConfig.Host = d.config.Host
	pgConnConfig.User = d.config.User
	pgConnConfig.Port = d.config.Port
	pgConnConfig.Password = d.config.Password
	pgConnConfig.Database = d.config.Database

	pgConn, err := pgx.ConnectConfig(context.Background(), pgConnConfig)
	if err != nil {
		return nil, err
	}

	return pgConn, nil
}

var (
	_ thunder.SourceDriver = (*Driver)(nil)
)
