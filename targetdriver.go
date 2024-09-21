package thunder

import (
	"github.com/rs/zerolog"
	"os"
)

type TargetDriverInfo struct {
	ID  string                                 `json:"ID"`
	New func(config any) (TargetDriver, error) `json:"-"`

	Name   string        `json:"name"`
	Config DynamicConfig `json:"-"`

	// As inlined SVG
	Image string   `json:"image,omitempty"`
	Notes []string `json:"notes,omitempty"`
}

type TargetDriver interface {
	DriverInfo() TargetDriverInfo

	TestConfig() (string, error) // TODO USELESS REPLACE IN FAVOR OF STATS TO CHECK NOT EMPTY

	HandleEvents(processor *Processor, eventsChan <-chan TargetEvent) error

	Shutdown() error
}

// Event

type TargetInsertEvent struct {
	Pkey string
	Json []byte
}

type TargetPatchEvent struct {
	Relation  *Relation
	Pkey      string
	JsonPatch []byte
}

type TargetDeleteEvent struct {
	Relation *Relation
	Pkey     string
}

type TargetTruncateEvent struct {
	Relation *Relation
}

type TargetEvent any // TargetDeleteEvent | TargetInsertEvent | TargetPatchEvent | TargetTruncateEvent

func GetLoggerForTargetDriver(t TargetDriver) *zerolog.Logger {
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel).With().Str("target-driver", t.DriverInfo().ID).Stack().Timestamp().Logger()
	return &logger
}
