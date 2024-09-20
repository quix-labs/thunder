package postgresql_flash

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/quix-labs/flash"
	"github.com/quix-labs/flash/drivers/wal_logical"
	_ "github.com/quix-labs/flash/drivers/wal_logical"
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

func (d *Driver) Start(p *thunder.Processor, eventsChan chan<- thunder.DbEvent, ctx context.Context) error {

	publicationSlotPrefix := "thunder_p" + strconv.Itoa(p.ID)
	replicationSlot := "thunder_replication_p" + strconv.Itoa(p.ID)

	conn, err := d.newConn()
	if err != nil {
		return err
	}

	// DROP PREVIOUS REPLICATION SLOT (fix flash bad closing)
	if _, err := conn.Exec(context.Background(), fmt.Sprintf("select pg_drop_replication_slot(slot_name) from pg_replication_slots where slot_name = '%s'", replicationSlot)); err != nil {
		return err
	}
	// DROP PREVIOUS PUBLICATION SLOT (fix flash bad closing)
	query := fmt.Sprintf(`
	DO $$
		DECLARE
			r RECORD;
		BEGIN
			FOR r IN (SELECT pubname FROM pg_publication WHERE pubname LIKE '%s%%') LOOP
				EXECUTE 'DROP PUBLICATION ' || quote_ident(r.pubname);
			END LOOP;
		END $$;`, publicationSlotPrefix)

	if _, err := conn.Exec(context.Background(), query); err != nil {
		return err
	}
	if err := conn.Close(context.Background()); err != nil {
		return err
	}

	// START LISTENING
	flashClient, _ := flash.NewClient(&flash.ClientConfig{
		DatabaseCnx: "postgresql://devuser:devpass@localhost:5432/devdb",
		Driver: wal_logical.NewDriver(&wal_logical.DriverConfig{
			PublicationSlotPrefix: publicationSlotPrefix,
			ReplicationSlot:       replicationSlot,
		}),
		//Driver: trigger.NewDriver(&trigger.DriverConfig{
		//	Schema: publicationSlotPrefix,
		//}),
	})

	// Register all listener
	prefixedConfig, err := GetRealtimeConfigForProcessor(p)
	if err != nil {
		return err
	}

	var listenerErrChan = make(chan error, 1)

	for path, config := range prefixedConfig {
		listener, err := flash.NewListener(config.ListenerConfig)
		if err != nil {
			return err
		}
		// Handle root table changes
		if path == "" {
			// TODO OFF GRACEFULL SHUTDOWN
			_, err = listener.On(flash.OperationAll, func(event flash.Event) {
				switch typedEvent := event.(type) {
				case *flash.InsertEvent:
					pkey, err := ExtractKeysFromMapAsJsonString(config.PrimaryKeys, *typedEvent.New)
					if err != nil {
						listenerErrChan <- err
						return
					}
					eventsChan <- &thunder.DbInsertEvent{Pkey: pkey}
				case *flash.UpdateEvent:
					pkey, err := ExtractKeysFromMapAsJsonString(config.PrimaryKeys, *typedEvent.Old)
					if err != nil {
						listenerErrChan <- err
						return
					}
					jsonPatch, err := json.Marshal(MapDiff(*typedEvent.Old, *typedEvent.New))
					if err != nil {
						listenerErrChan <- err
						return
					}
					eventsChan <- &thunder.DbPatchEvent{Pkey: pkey, JsonPatch: jsonPatch}
				case *flash.DeleteEvent:
					pkey, err := ExtractKeysFromMapAsJsonString(config.PrimaryKeys, *typedEvent.Old)
					if err != nil {
						listenerErrChan <- err
						return
					}
					eventsChan <- &thunder.DbDeleteEvent{Pkey: pkey}
				case *flash.TruncateEvent:
					eventsChan <- &thunder.DbTruncateEvent{}
				}
			})
			if err != nil {
				return err
			}
			flashClient.Attach(listener)
			continue
		}

		// HANDLE RELATION EVENTS
		_, err = listener.On(flash.OperationUpdate^flash.OperationTruncate^flash.OperationDelete, func(event flash.Event) {
			switch typedEvent := event.(type) {
			case *flash.UpdateEvent:
				pkey, err := ExtractKeysFromMapAsJsonString(config.PrimaryKeys, *typedEvent.Old)
				if err != nil {
					listenerErrChan <- err
					return
				}
				jsonPatch, err := json.Marshal(MapDiff(*typedEvent.Old, *typedEvent.New))
				if err != nil {
					listenerErrChan <- err
					return
				}
				eventsChan <- &thunder.DbPatchEvent{
					Path:      path,
					Pkey:      pkey,
					JsonPatch: jsonPatch,
				}

			case *flash.TruncateEvent:
				eventsChan <- &thunder.DbTruncateEvent{
					Path: path,
				}

			case *flash.DeleteEvent:
				pkey, err := ExtractKeysFromMapAsJsonString(config.PrimaryKeys, *typedEvent.Old)
				if err != nil {
					listenerErrChan <- err
					return
				}

				eventsChan <- &thunder.DbDeleteEvent{
					Path: path,
					Pkey: pkey,
				}
			}

		})
		if err != nil {
			return err
		}
		flashClient.Attach(listener)
	}

	// START IN BACKGROUND
	errChan := make(chan error)
	go func() {
		errChan <- flashClient.Start()
	}()

	// Start listening
	for {
		select {
		case <-ctx.Done():
			if err := flashClient.Close(); err != nil {
				return err
			}
			return ctx.Err()
		case err := <-errChan:
			return err
		case err := <-listenerErrChan:
			return err
		}
	}
}

func (d *Driver) Stop() error {
	return nil
}

func (d *Driver) newConn() (*pgx.Conn, error) {
	pgConnConfig, err := pgx.ParseConfig("postgres://u:s@l:5432/d?sslmode=disable") // Keep sslMode to work in scratch docker
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
