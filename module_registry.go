package thunder

import (
	"errors"
	"fmt"
	"sync"
)

var (
	modules   = make(map[string]ModuleInfo)
	modulesMu sync.RWMutex
)

func RegisterModule(module Module) {
	info := module.ThunderModule()
	if info.ID == "" {
		panic("module ID missing")
	}
	if info.New == nil {
		panic("module New function missing")
	}
	if val := info.New(); val == nil {
		panic("ModuleInfo.New must return a non-nil module instance")
	}

	modulesMu.Lock()
	defer modulesMu.Unlock()
	if _, ok := modules[info.ID]; ok {
		panic(fmt.Sprintf("module already registered: %s", info.ID))
	}
	modules[info.ID] = info
}

func GetModules() []ModuleInfo {
	modulesMu.RLock()
	defer modulesMu.RUnlock()
	var mods []ModuleInfo
	for _, info := range modules {
		mods = append(mods, info)
	}
	return mods
}

func GetModule(ID string) (ModuleInfo, error) {
	modulesMu.RLock()
	defer modulesMu.RUnlock()

	info, ok := modules[ID]
	if !ok {
		return ModuleInfo{}, errors.New("module not registered: " + ID)
	}
	return info, nil
}
