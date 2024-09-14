package thunder

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

var (
	configMu sync.RWMutex
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

// JSON TRANSFORMATION

type JsonConfig struct {
	Sources    []*JsonSource    `json:"sources"`
	Targets    []*JsonTarget    `json:"targets"`
	Processors []*JsonProcessor `json:"processors"`
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

	var jsonConfig JsonConfig
	if err = json.Unmarshal(content, &jsonConfig); err != nil {
		return err
	}

	// Parse jsonSource config
	for _, jsonSource := range jsonConfig.Sources {
		source, err := UnserializeSource(jsonSource)
		if err != nil {
			return err
		}
		if err = AddSource(source); err != nil {
			return err
		}
	}

	// Parse jsonTargets config
	for _, jsonTarget := range jsonConfig.Targets {
		target, err := UnserializeTarget(jsonTarget)
		if err != nil {
			return err
		}

		if err = AddTarget(target); err != nil {
			return err
		}
	}

	// Parse jsonProcessors config
	for _, jsonProcessor := range jsonConfig.Processors {
		processor, err := UnserializeProcessor(jsonProcessor)
		if err != nil {
			return err
		}

		if err = AddProcessor(processor); err != nil {
			return err
		}
	}

	return nil
}

func SaveConfig() error {
	var config JsonConfig

	// Add sources
	config.Sources = []*JsonSource{}
	for _, source := range GetSources() {
		jsonSource, err := SerializeSource(source)
		if err != nil {
			return err
		}
		config.Sources = append(config.Sources, jsonSource)
	}

	// Add Processors
	config.Processors = []*JsonProcessor{}
	for _, processor := range GetProcessors() {
		jsonProcessor, err := SerializeProcessor(processor)
		if err != nil {
			return err
		}
		config.Processors = append(config.Processors, jsonProcessor)
	}

	// Add targets
	config.Targets = []*JsonTarget{}
	for _, target := range GetTargets() {
		jsonTarget, err := SerializeTarget(target)
		if err != nil {
			return err
		}
		config.Targets = append(config.Targets, jsonTarget)
	}

	bytes, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}

	if err = os.WriteFile(GetConfigPath(), bytes, 0644); err != nil {
		return err
	}
	return nil
}
