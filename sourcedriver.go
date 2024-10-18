package thunder

import (
	"encoding/json"
	"github.com/creasty/defaults"
	"github.com/quix-labs/thunder/utils"
	"github.com/rs/zerolog"
	"os"
	"reflect"
)

type SourceDriverInfo struct {
	ID  string                                 `json:"ID"`
	New func(config any) (SourceDriver, error) `json:"-"`

	Name   string        `json:"name"`
	Config DynamicConfig `json:"-"`

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
	DriverInfo() SourceDriverInfo

	TestConfig() (string, error) // TODO USELESS REPLACE IN FAVOR OF STATS TO CHECK NOT EMPTY
	Stats() (*SourceDriverStats, error)

	// GetDocumentsForProcessor use limit=0 to unlimited
	GetDocumentsForProcessor(processor *Processor, in utils.BroadcasterIn[*Document], limit uint64) error

	// Real Time Indexing

	Start(processor *Processor, in utils.BroadcasterIn[DbEvent]) error
	Stop() error
}

// Events

type DbInsertEvent struct {
	Pkey string
}

type DbPatchEvent struct {
	Relation *Relation

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

func ConvertSourceDriverConfig(driver *SourceDriverInfo, config any) (any, error) {
	typedConfig := reflect.New(reflect.TypeOf((*driver).Config)).Interface()
	bytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(bytes, typedConfig); err != nil {
		return nil, err
	}

	// Apply default tag if needed
	if err := defaults.Set(typedConfig); err != nil {
		return nil, err
	}

	return typedConfig, nil
}

func GetLoggerForSourceDriver(s SourceDriver) *zerolog.Logger {
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel).With().Str("source-driver", s.DriverInfo().ID).Stack().Timestamp().Logger()
	return &logger
}
