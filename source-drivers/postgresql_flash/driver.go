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

func (d *Driver) ThunderSourceDriver() thunder.SourceDriverInfo {
	return thunder.SourceDriverInfo{
		ID:  "postgresql_flash",
		New: func() thunder.SourceDriver { return new(Driver) },

		Name:   "PostgreSQL (Flash)",
		Image:  logo,
		Config: DriverConfig{},
		Notes:  []string{"Please be sure wal_replication is set to logical on your PostgreSQL config"},
	}
}

func (d *Driver) TestConfig(config any) (string, error) {
	cfg, ok := config.(*DriverConfig)
	if !ok {
		return "", errors.New("invalid config type")
	}

	pgConn, err := d.GetConn(cfg)
	if err != nil {
		return "Cannot connect to database", err
	}
	defer pgConn.Close(context.Background())

	return "Successfully connected", nil
}

func (d *Driver) Stats(config any) (*thunder.SourceDriverStats, error) {
	cfg, ok := config.(*DriverConfig)
	if !ok {
		return nil, errors.New("invalid config type")
	}
	pgConn, err := d.GetConn(cfg)
	if err != nil {
		return nil, err
	}
	defer pgConn.Close(context.Background())

	query := `SELECT
    table_name as name,
    JSON_AGG(column_name) AS columns,
    TO_JSON(ARRAY_AGG(column_name) FILTER (WHERE primary_key)) AS primary_keys
FROM (
         SELECT
             pgc.relname AS table_name,
             a.attname AS column_name,
             COALESCE(i.indisprimary, false) AS primary_key
         FROM
             pg_attribute a
                 JOIN pg_class pgc ON pgc.oid = a.attrelid
                 LEFT JOIN pg_index i ON (pgc.oid = i.indrelid AND a.attnum = ANY(i.indkey))
                 LEFT JOIN pg_catalog.pg_namespace n ON n.oid = pgc.relnamespace
         WHERE
        	pgc.relkind IN ('r', '')  -- Relkind for tables
            AND n.nspname <> 'pg_catalog'
        	AND n.nspname <> 'information_schema'
           	AND n.nspname !~ '^pg_toast'
           	AND a.attnum > 0
           	AND pg_table_is_visible(pgc.oid)
           	AND NOT a.attisdropped
     ) AS subquery
GROUP BY table_name
ORDER BY table_name`

	stats := thunder.SourceDriverStats{}
	type RowResult struct {
		Name        string   `json:"name"`
		Columns     []string `json:"columns"`
		PrimaryKeys []string `json:"primary_keys"`
	}

	results, err := GetResultsSync[RowResult](pgConn, query, time.Second*10)
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

func (d *Driver) Start(config any) error {
	var ok bool
	if d.config, ok = config.(*DriverConfig); !ok {
		return errors.New("invalid config type")
	}
	var err error
	if d.conn, err = d.GetConn(d.config); err != nil {
		return err
	}
	return nil
}

func (d *Driver) GetDocumentsForProcessor(processor *thunder.Processor, limit uint64) (<-chan *thunder.Document, <-chan error) {
	query, err := GetSqlForMapping(processor.Table, &processor.Mapping)
	if err != nil {
		panic(err)
	}
	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT %s", query, strconv.FormatUint(limit, 10))
	}
	return GetResultsChan[thunder.Document](d.conn, query)
}

func (d *Driver) Stop() error {
	if d.conn != nil {
		return d.conn.Close(context.Background())
	}
	return nil
}

func (d *Driver) GetConn(cfg *DriverConfig) (*pgx.Conn, error) {
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
