package thunder

import (
	"errors"
	"fmt"
	"sync"
)

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
