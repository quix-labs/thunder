package thunder

import (
	"context"
	"encoding/json"
	"github.com/creasty/defaults"
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

	GetDocumentsForProcessor(processor *Processor, docChan chan<- *Document, errChan chan error, limit uint64)

	// Real Time Indexing

	Start(processor *Processor, eventsChan chan<- DbEvent, ctx context.Context) error
	Stop() error
}

// Events

type DbInsertEvent struct {
	Pkey string
}

type DbPatchEvent struct {
	Path string

	Pkey      string
	JsonPatch []byte
}

type DbDeleteEvent struct {
	Pkey string
}

type DbTruncateEvent struct{}

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
