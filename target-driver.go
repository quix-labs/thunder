package thunder

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/creasty/defaults"
	"reflect"
	"sync"
)

type TargetDriverInfo struct {
	ID  string                                 `json:"ID"`
	New func(config any) (TargetDriver, error) `json:"-"`

	Name   string `json:"name"`
	Config any    `json:"-"`
	// As inlined SVG
	Image string   `json:"image,omitempty"`
	Notes []string `json:"notes,omitempty"`
}

type TargetDriver interface {
	DriverInfo() TargetDriverInfo

	TestConfig() (string, error) // TODO USELESS REPLACE IN FAVOR OF STATS TO CHECK NOT EMPTY
}

// UTILITIES FUNCTIONS

func ConvertTargetDriverConfig(driver *TargetDriverInfo, config any) (any, error) {
	typedConfig := reflect.New(reflect.TypeOf((*driver).Config)).Interface()
	bytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(bytes, &typedConfig); err != nil {
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
	targetDrivers   = make(map[string]TargetDriverInfo)
	targetDriversMu sync.RWMutex
)

func RegisterTargetDriver(driver TargetDriver) {
	info := driver.DriverInfo()
	if info.ID == "" {
		panic("target driver ID missing")
	}
	if info.New == nil {
		panic("target driver New function missing")
	}

	targetDriversMu.Lock()
	defer targetDriversMu.Unlock()
	if _, ok := targetDrivers[info.ID]; ok {
		panic(fmt.Sprintf("target driver already registered: %s", info.ID))
	}
	targetDrivers[info.ID] = info
}

func GetTargetDrivers() []TargetDriverInfo {
	targetDriversMu.RLock()
	defer targetDriversMu.RUnlock()
	var infos []TargetDriverInfo
	for _, info := range targetDrivers {
		infos = append(infos, info)
	}
	return infos
}

func GetTargetDriver(ID string) (TargetDriverInfo, error) {
	targetDriversMu.RLock()
	defer targetDriversMu.RUnlock()

	info, ok := targetDrivers[ID]
	if !ok {
		return TargetDriverInfo{}, errors.New("target driver not registered: " + ID)
	}
	return info, nil
}
