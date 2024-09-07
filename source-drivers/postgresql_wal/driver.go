package postgresql_wal

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/quix-labs/thunder"
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
type Driver struct{}

//go:embed icon.svg
var logo string

func (d *Driver) ThunderSourceDriver() thunder.SourceDriverInfo {
	return thunder.SourceDriverInfo{
		ID:  "postgresql_wal",
		New: func() thunder.SourceDriver { return new(Driver) },

		Name:   "PostgreSQL (WAL)",
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

	rows := pgConn.ExecParams(context.Background(), `SELECT
    table_name as name,
    JSON_AGG(column_name)::TEXT AS columns,
    TO_JSON(ARRAY_AGG(column_name) FILTER (WHERE primary_key))::TEXT AS primarykey
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
ORDER BY table_name;`, nil, nil, nil, nil)

	stats := thunder.SourceDriverStats{}
	for rows.NextRow() {
		tableName := string(rows.Values()[0])
		var columns []string
		var primaryKeys []string

		_ = json.Unmarshal(rows.Values()[1], &columns)
		_ = json.Unmarshal(rows.Values()[2], &primaryKeys)

		stats[tableName] = thunder.SourceDriverStatsTable{
			Columns:     columns,
			PrimaryKeys: primaryKeys,
		}
	}

	_, err = rows.Close()
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (d *Driver) GetConn(cfg *DriverConfig) (*pgconn.PgConn, error) {
	pgConnConfig, err := pgconn.ParseConfig("")
	if err != nil {
		return nil, err
	}
	pgConnConfig.Host = cfg.Host
	pgConnConfig.User = cfg.User
	pgConnConfig.Port = cfg.Port
	pgConnConfig.Password = cfg.Password
	pgConnConfig.Database = cfg.Database

	pgConn, err := pgconn.ConnectConfig(context.Background(), pgConnConfig)
	if err != nil {
		return nil, err
	}
	//defer pgConn.Close(context.Background()) // TODO OUTSIDE OR RETURN CLOSE
	return pgConn, nil
}

var (
	_ thunder.SourceDriver = (*Driver)(nil) // Interface implementation
)
