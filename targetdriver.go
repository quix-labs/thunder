package thunder

import (
	"errors"
	"github.com/quix-labs/thunder/utils"
	"github.com/rs/zerolog"
	"os"
)

type TargetDriverConfig struct {
	Name   string              `json:"name"`
	Config utils.DynamicConfig `json:"-"`

	// As inlined SVG
	Image string   `json:"image,omitempty"`
	Notes []string `json:"notes,omitempty"`
}

type TargetDriver interface {
	ID() string

	New(config any) (TargetDriver, error)

	Config() TargetDriverConfig

	TestConfig() (string, error) // TODO USELESS REPLACE IN FAVOR OF STATS TO CHECK NOT EMPTY

	HandleEvents(processor *Processor, eventsChan <-chan TargetEvent) error

	Shutdown() error
}

// Event

type TargetInsertEvent struct {
	Pkey    string
	Version int
	Json    []byte
}

type TargetPatchEvent struct {
	Relation  *Relation
	Pkey      string
	Version   int
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

// UTILITIES FUNCTIONS

func GetLoggerForTargetDriver(driver TargetDriver) *zerolog.Logger {
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel).With().Str("target-driver", driver.ID()).Stack().Timestamp().Logger()
	return &logger
}

// TargetDrivers is a registry that allows external library to register their own target driver
var TargetDrivers = utils.NewRegistry[TargetDriver]("target-driver").ValidateUsing(func(ID string, driver TargetDriver) error {
	if driver.New == nil {
		return errors.New(`target driver doesn't implements "New" method`)
	}
	return nil
})
