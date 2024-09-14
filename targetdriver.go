package thunder

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

	IndexDocumentsForProcessor(processor *Processor, docChan <-chan *Document, errChan chan error)

	Shutdown() error
}
