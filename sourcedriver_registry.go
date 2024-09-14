package thunder

import (
	"errors"
	"fmt"
	"sync"
)

// REGISTERING MODULES
var (
	sourceDrivers   = make(map[string]SourceDriverInfo)
	sourceDriversMu sync.RWMutex
)

func RegisterSourceDriver(sourceDriver SourceDriver) {
	info := sourceDriver.DriverInfo()
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
