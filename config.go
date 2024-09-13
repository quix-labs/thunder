package thunder

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type Source struct {
	Driver string        `json:"driver"`
	Config DynamicConfig `json:"config"`
}

type Target struct {
	Driver string        `json:"driver"`
	Config DynamicConfig `json:"config"`
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

	var jsonConfig struct {
		Sources []struct {
			Driver string `json:"driver"`
			Config any    `json:"config"`
		} `json:"sources"`
		Processors []Processor `json:"processors"`
		Targets    []struct {
			Driver string `json:"driver"`
			Config any    `json:"config"`
		} `json:"targets"`
	}

	if err = json.Unmarshal(content, &jsonConfig); err != nil {
		return err
	}

	// Parse source config
	config.Sources = make([]Source, len(jsonConfig.Sources))
	for idx, source := range jsonConfig.Sources {
		if driver, err := GetSourceDriver(source.Driver); err == nil {
			if typedConfig, err := ConvertToDynamicConfig(&driver.Config, source.Config); err == nil {
				config.Sources[idx] = Source{
					Driver: source.Driver,
					Config: typedConfig,
				}
			}
		}
	}

	// Assign processors
	config.Processors = jsonConfig.Processors

	// Parse target config
	config.Targets = make([]Target, len(jsonConfig.Targets))
	for idx, target := range jsonConfig.Targets {
		if driver, err := GetTargetDriver(target.Driver); err == nil {
			if typedConfig, err := ConvertToDynamicConfig(&driver.Config, target.Config); err == nil {
				config.Targets[idx] = Target{
					Driver: target.Driver,
					Config: typedConfig,
				}
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
