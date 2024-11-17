package mysql

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/utils"
	"strconv"
	"strings"
)

const DriverID = "thunder.mysql"

func init() {
	_ = thunder.SourceDrivers.Register(DriverID, &Driver{})
}

type Driver struct {
	config *DriverConfig
}

func (d *Driver) New(config any) (thunder.SourceDriver, error) {
	cfg, ok := config.(*DriverConfig)
	if !ok {
		return nil, errors.New("invalid config type")
	}
	return &Driver{config: cfg}, nil
}

func (d *Driver) ID() string {
	return DriverID
}

//go:embed icon.svg
var logo string

func (d *Driver) Config() thunder.SourceDriverConfig {
	return thunder.SourceDriverConfig{
		Name:   "MySQL",
		Image:  logo,
		Config: DriverConfig{},
		Notes:  []string{"This driver is not intended to support real-time sync"},
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
	conn, err := d.newConn(context.Background())
	if err != nil {
		return nil, err
	}

	query := StatsQuery(d.config.Database)
	rows, err := conn.QueryxContext(context.Background(), query)
	if err != nil {
		return nil, err
	}

	type RowResult struct {
		Name        string `db:"name"`
		Columns     string `db:"columns"`
		PrimaryKeys string `db:"primary_keys"`
	}

	var results []*RowResult
	err = sqlx.StructScan(rows, &results)
	if err != nil {
		return nil, err
	}

	closeErr := conn.Close()

	stats := thunder.SourceDriverStats{}
	for _, result := range results {
		stats[result.Name] = thunder.SourceDriverStatsTable{
			Columns:     strings.Split(result.Columns, ","),
			PrimaryKeys: strings.Split(result.PrimaryKeys, ","),
		}
	}

	return &stats, closeErr
}

func (d *Driver) GetDocumentsForProcessor(processor *thunder.Processor, in chan<- *thunder.Document, ctx context.Context, limit uint64) error {
	query, err := GetSqlForProcessor(processor)
	if err != nil {
		return err
	}
	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT %s", query, strconv.FormatUint(limit, 10))
	}

	conn, err := d.newConn(ctx)
	if err != nil {
		return err
	}

	rows, err := conn.QueryxContext(ctx, query)
	if err != nil {
		return err
	}

	for rows.Next() {
		var rowResult = thunder.Document{}
		if err = rows.StructScan(&rowResult); err != nil {
			break
		}
		in <- &rowResult
	}

	return errors.Join(err, rows.Close(), conn.Close())
}

func (d *Driver) Start(p *thunder.Processor, in utils.BroadcasterIn[thunder.DbEvent]) error {
	return errors.New("not implemented")
}

func (d *Driver) Stop() error {
	return errors.New("not implemented")
}

func (d *Driver) newConn(ctx context.Context) (*sqlx.Conn, error) {

	cfg := mysql.Config{
		User:   d.config.User,
		Passwd: d.config.Password,
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%d", d.config.Host, d.config.Port),
		DBName: d.config.Database,
	}

	db, err := sqlx.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	return db.Connx(ctx)
}

var (
	_ thunder.SourceDriver = (*Driver)(nil)
)
