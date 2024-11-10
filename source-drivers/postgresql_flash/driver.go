package postgresql_flash

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/quix-labs/flash"
	"github.com/quix-labs/flash/drivers/trigger"
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/utils"
	"strconv"
	"time"
)

const DriverID = "thunder.postgresql_flash"

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
	conn, err := d.newConn(context.TODO())
	if err != nil {
		return nil, err
	}

	query := StatsQuery(d.config.Schema)
	stats := thunder.SourceDriverStats{}
	type RowResult struct {
		Name        string   `json:"name"`
		Columns     []string `json:"columns"`
		PrimaryKeys []string `json:"primary_keys"`
	}

	results, err := GetResultsSync[RowResult](conn, query, time.Second*10, false)
	closeErr := conn.Close(context.Background())

	if err != nil {
		return nil, errors.Join(err, closeErr)
	}
	for _, result := range results {
		stats[result.Name] = thunder.SourceDriverStatsTable{
			Columns:     result.Columns,
			PrimaryKeys: result.PrimaryKeys,
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

	resultErr := GetResult[thunder.Document](conn, query, in, ctx)
	closeErr := conn.Close(ctx)
	return errors.Join(resultErr, closeErr)
}

func (d *Driver) Start(p *thunder.Processor, in utils.BroadcasterIn[thunder.DbEvent]) error {

	publicationSlotPrefix := "thunder_p" + p.ID
	replicationSlot := "thunder_replication_p" + p.ID

	conn, err := d.newConn(context.TODO())
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
	//time.Sleep(time.Second * 10) // TODO FIX BUG IN FLASH

	// START LISTENING
	flashClient, _ := flash.NewClient(&flash.ClientConfig{
		Logger:      thunder.GetLoggerForSourceDriver(DriverID),
		DatabaseCnx: "postgresql://devuser:devpass@localhost:5432/devdb?sslmode=disable",
		//Driver: wal_logical.NewDriver(&wal_logical.DriverConfig{
		//	PublicationSlotPrefix: publicationSlotPrefix,
		//	ReplicationSlot:       replicationSlot,
		//}),
		Driver: trigger.NewDriver(&trigger.DriverConfig{
			Schema: fmt.Sprintf("thunder%s", p.ID),
		}),
	})

	// Register all listener
	mappedConfig, err := GetRealtimeConfigForProcessor(p)
	if err != nil {
		return err
	}

	var listenerErrChan = make(chan error, 1)

	for relation, config := range mappedConfig {
		listener, err := flash.NewListener(config.ListenerConfig)
		if err != nil {
			return err
		}

		// Handle root table changes
		if relation == nil {
			// TODO OFF GRACEFULL SHUTDOWN
			_, err = listener.On(flash.OperationAll, func(event flash.Event) {
				switch typedEvent := event.(type) {
				case *flash.InsertEvent:
					pkey, err := ExtractPkeyFromMap(config.PrimaryKeys, *typedEvent.New)
					if err != nil {
						listenerErrChan <- err
						return
					}
					in.Broadcast(&thunder.DbInsertEvent{Pkey: pkey})
				case *flash.UpdateEvent:
					pkey, err := ExtractPkeyFromMap(config.PrimaryKeys, *typedEvent.Old)
					if err != nil {
						listenerErrChan <- err
						return
					}
					jsonPatch, err := json.Marshal(MapDiff(*typedEvent.Old, *typedEvent.New))
					if err != nil {
						listenerErrChan <- err
						return
					}
					in.Broadcast(&thunder.DbPatchEvent{Pkey: pkey, JsonPatch: jsonPatch})
				case *flash.DeleteEvent:
					pkey, err := ExtractPkeyFromMap(config.PrimaryKeys, *typedEvent.Old)
					if err != nil {
						listenerErrChan <- err
						return
					}
					in.Broadcast(&thunder.DbDeleteEvent{Pkey: pkey})
				case *flash.TruncateEvent:
					in.Broadcast(&thunder.DbTruncateEvent{})
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
				pkey, err := ExtractPkeyFromMap(config.PrimaryKeys, *typedEvent.Old)
				if err != nil {
					listenerErrChan <- err
					return
				}
				jsonPatch, err := json.Marshal(MapDiff(*typedEvent.Old, *typedEvent.New))
				if err != nil {
					listenerErrChan <- err
					return
				}

				in.Broadcast(&thunder.DbPatchEvent{
					Relation:  relation,
					Pkey:      pkey,
					JsonPatch: jsonPatch,
				})

			case *flash.TruncateEvent:
				in.Broadcast(&thunder.DbTruncateEvent{
					Relation: relation,
				})

			case *flash.DeleteEvent:
				pkey, err := ExtractPkeyFromMap(config.PrimaryKeys, *typedEvent.Old)
				if err != nil {
					listenerErrChan <- err
					return
				}

				in.Broadcast(&thunder.DbDeleteEvent{
					Relation: relation,
					Pkey:     pkey,
				})
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
	defer func() {
		thunder.GetLoggerForSourceDriver(DriverID).Debug().Msg("send close signal to flash")
		if err := flashClient.Close(); err != nil {
			thunder.GetLoggerForSourceDriver(DriverID).Error().Msg(err.Error())
		}
	}()

	// Start listening
	for {
		select {
		case err := <-errChan:
			return err
		case err := <-listenerErrChan:
			return err
		default:
			if in.Closed() {
				return nil
			}
		}
	}
}

func (d *Driver) Stop() error {
	return nil
}

func (d *Driver) newConn(ctx context.Context) (*pgx.Conn, error) {
	pgConnConfig, err := pgx.ParseConfig("postgres://u:s@l:5432/d?sslmode=disable") // Keep sslMode to work in scratch docker
	if err != nil {
		return nil, err
	}
	pgConnConfig.Host = d.config.Host
	pgConnConfig.User = d.config.User
	pgConnConfig.Port = d.config.Port
	pgConnConfig.Password = d.config.Password
	pgConnConfig.Database = d.config.Database

	pgConn, err := pgx.ConnectConfig(ctx, pgConnConfig)
	if err != nil {
		return nil, err
	}

	return pgConn, nil
}

var (
	_ thunder.SourceDriver = (*Driver)(nil)
)
