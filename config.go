package thunder

import (
	"encoding/json"
	"os"
	"sync"
)

type Source struct {
	Driver string `json:"driver"`
	Config any    `json:"config"`
}

type Processor struct {
	Source int    `json:"source"`
	Table  string `json:"table"`

	Mapping []any `json:"mapping"`

	Index string `json:"index"`
}

type Config struct {
	Sources    []Source    `json:"sources"`
	Processors []Processor `json:"processors"`
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
	content, err := os.ReadFile(GetConfigPath())
	if err != nil {
		return err
	}

	configMu.Lock()
	defer configMu.Unlock()
	if err = json.Unmarshal(content, &config); err != nil {
		return err
	}

	// Parse driver config
	for idx, source := range config.Sources {
		if driver, err := GetSourceDriver(source.Driver); err == nil {
			if typedConfig, err := ConvertSourceDriverConfig(&driver, source.Config); err == nil {
				config.Sources[idx].Config = typedConfig
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
