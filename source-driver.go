package thunder

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/creasty/defaults"
	"reflect"
	"sync"
)

type SourceDriverInfo struct {
	ID  string                                 `json:"ID"`
	New func(config any) (SourceDriver, error) `json:"-"`

	Name   string `json:"name"`
	Config any    `json:"-"`
	// As inlined SVG
	Image string   `json:"image,omitempty"`
	Notes []string `json:"notes,omitempty"`
}

type SourceDriverStatsTable struct {
	Columns     []string `json:"columns"`
	PrimaryKeys []string `json:"primary_keys"`
}

type SourceDriverStats map[string]SourceDriverStatsTable
type Document struct {
	Key  string `json:"id"`
	Json []byte `json:"json"`
}

type SourceDriver interface {
	ThunderSourceDriver() SourceDriverInfo

	TestConfig() (string, error) // TODO USELESS REPLACE IN FAVOR OF STATS TO CHECK NOT EMPTY
	Stats() (*SourceDriverStats, error)

	GetDocumentsForProcessor(processor *Processor, limit uint64) (<-chan *Document, <-chan error)

	Start() error
	Stop() error

	Shutdown() error
}

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

// REGISTERING MODULES
var (
	sourceDrivers   = make(map[string]SourceDriverInfo)
	sourceDriversMu sync.RWMutex
)

func RegisterSourceDriver(sourceDriver SourceDriver) {
	info := sourceDriver.ThunderSourceDriver()
	if info.ID == "" {
		panic("source driver ID missing")
	}
	if info.New == nil {
		panic("source driver New function missing")
	}

	sourceDriversMu.Lock()
	defer sourceDriversMu.Unlock()
	if _, ok := sourceDrivers[info.ID]; ok {
		panic(fmt.Sprintf("source driver already registered: %s", info.ID))
	}
	sourceDrivers[info.ID] = info
}

func GetSourceDrivers() []SourceDriverInfo {
	sourceDriversMu.RLock()
	defer sourceDriversMu.RUnlock()
	var infos []SourceDriverInfo
	for _, info := range sourceDrivers {
		infos = append(infos, info)
	}
	return infos
}

func GetSourceDriver(ID string) (SourceDriverInfo, error) {
	sourceDriversMu.RLock()
	defer sourceDriversMu.RUnlock()

	info, ok := sourceDrivers[ID]
	if !ok {
		return SourceDriverInfo{}, errors.New("source driver not registered: " + ID)
	}
	return info, nil
}
