package thunder

import "context"

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

	HandleEvents(processor *Processor, eventsChan <-chan TargetEvent, ctx context.Context) error

	Shutdown() error
}

// Event

type TargetInsertEvent struct {
	Pkey string
	Json []byte
}

type TargetPatchEvent struct {
	Path      string
	Pkey      string
	JsonPatch []byte
}

type TargetDeleteEvent struct {
	Path string
	Pkey string
}

type TargetTruncateEvent struct {
	Path string
}

type TargetEvent any // TargetDeleteEvent | TargetInsertEvent | TargetPatchEvent | TargetTruncateEvent
