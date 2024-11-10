package thunder

import (
	"context"
	"errors"
	"github.com/quix-labs/thunder/utils"
	"github.com/rs/zerolog"
	"os"
)

type SourceDriverConfig struct {
	Name   string              `json:"name"`
	Config utils.DynamicConfig `json:"-"`

	// As inlined SVG
	Image string   `json:"image,omitempty"`
	Notes []string `json:"notes,omitempty"`
}

type SourceDriverStatsTable struct {
	Columns     []string `json:"columns"`
	PrimaryKeys []string `json:"primary_keys"`
}

type SourceDriverStats map[string]SourceDriverStatsTable

type SourceDriver interface {
	ID() string

	New(config any) (SourceDriver, error)

	Config() SourceDriverConfig

	TestConfig() (string, error) // TODO USELESS REPLACE IN FAVOR OF STATS TO CHECK NOT EMPTY
	Stats() (*SourceDriverStats, error)

	// GetDocumentsForProcessor use limit=0 to unlimited
	GetDocumentsForProcessor(processor *Processor, in chan<- *Document, ctx context.Context, limit uint64) error

	// Real Time Indexing

	Start(processor *Processor, in utils.BroadcasterIn[DbEvent]) error
	Stop() error
}

// Events

type DbInsertEvent struct {
	Pkey    string
	Version int
}

type DbPatchEvent struct {
	Relation  *Relation
	Version   int
	Pkey      string
	JsonPatch []byte
}

type DbDeleteEvent struct {
	Relation *Relation
	Pkey     string
}

type DbTruncateEvent struct {
	Relation *Relation
}

type DbEvent any // DbDeleteEvent | DbInsertEvent | DbPatchEvent | DbTruncateEvent

// UTILITIES FUNCTIONS

func GetLoggerForSourceDriver(driverID string) *zerolog.Logger {
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel).With().Str("source-driver", driverID).Stack().Timestamp().Logger()
	return &logger
}

// SourceDrivers is a registry that allows external library to register their own source driver
var SourceDrivers = utils.NewRegistry[SourceDriver]("source-driver").ValidateUsing(func(ID string, driver SourceDriver) error {
	if driver.New == nil {
		return errors.New(`target driver doesn't implements "New" method`)
	}
	return nil
})
