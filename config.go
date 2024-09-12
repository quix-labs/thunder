package thunder

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type Source struct {
	Driver string `json:"driver"`
	Config any    `json:"config"`
}

type Target struct {
	Driver string `json:"driver"`
	Config any    `json:"config"`
}

type Config struct {
	Sources    []Source    `json:"sources"`
	Processors []Processor `json:"processors"`
	Targets    []Target    `json:"targets"`
}

var (
	configMu sync.RWMutex
	config   Config
	path     = "./config.json"
)

func GetConfigPath() string {
	configMu.RLock()
	defer configMu.RUnlock()
	return path
}

func SetConfigPath(newPath string) {
	configMu.RLock()
	defer configMu.RUnlock()
	path = newPath
}

func LoadConfig() error {
	if _, err := os.Stat(GetConfigPath()); errors.Is(err, os.ErrNotExist) {
		fmt.Println("config file doesn't exists, ignore load")
		return nil
	} else if err != nil {
		return err
	}

	content, err := os.ReadFile(GetConfigPath())
	if err != nil {
		return err
	}

	configMu.Lock()
	defer configMu.Unlock()
	if err = json.Unmarshal(content, &config); err != nil {
		return err
	}

	// Parse source config
	for idx, source := range config.Sources {
		if driver, err := GetSourceDriver(source.Driver); err == nil {
			if typedConfig, err := ConvertSourceDriverConfig(&driver, source.Config); err == nil {
				config.Sources[idx].Config = typedConfig
			}
		}
	}
	// Parse target config
	for idx, target := range config.Targets {
		if driver, err := GetTargetDriver(target.Driver); err == nil {
			if typedConfig, err := ConvertTargetDriverConfig(&driver, target.Config); err == nil {
				config.Targets[idx].Config = typedConfig
			}
		}
	}
	return nil
}

func SaveConfig() error {
	configMu.RLock()
	defer configMu.RUnlock()

	b, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}
	if err := os.WriteFile(GetConfigPath(), b, 0666); err != nil {
		return err
	}
	return nil
}

func GetConfig() Config {
	configMu.RLock()
	defer configMu.RUnlock()
	return config
}

func SetConfig(new Config) {
	configMu.Lock()
	defer configMu.Unlock()
	config = new
}
