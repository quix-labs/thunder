package utils

import (
	"cmp"
	"encoding/json"
	"github.com/creasty/defaults"
	"reflect"
	"strings"
)

type DynamicConfig interface {
	Excerpt() string
}
type DynamicConfigFields []DynamicConfigField
type DynamicConfigField struct {
	Name     string  `json:"name"`
	Label    string  `json:"label"`
	Type     string  `json:"type"`
	Required bool    `json:"required"`
	Help     *string `json:"help,omitempty"`
	Min      string  `json:"min,omitempty"`
	Max      string  `json:"max,omitempty"`
	Default  *string `json:"default,omitempty"`

	Options []string `json:"options,omitempty"`
}

func ParseDynamicConfigFields(config *DynamicConfig) DynamicConfigFields {
	var fields DynamicConfigFields
	configValue := reflect.ValueOf(*config)
	configType := reflect.TypeOf(*config)

	if configType.Kind() == reflect.Ptr {
		configValue = configValue.Elem()
		configType = configType.Elem()
	}

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)

		label := cmp.Or(field.Tag.Get("label"), field.Name)
		inputType := cmp.Or(field.Tag.Get("type"), "text")

		helpTag := field.Tag.Get("help")
		var helpText *string = nil
		if helpTag != "" {
			helpText = &helpTag
		}

		defaultTag := field.Tag.Get("default")
		var defaultValue *string = nil
		if defaultTag != "" {
			defaultValue = &defaultTag
		}

		var options []string
		if inputType == "select" {
			optionsTag := cmp.Or(field.Tag.Get("options"), defaultTag)
			if optionsTag != "" {
				options = strings.Split(optionsTag, ",")
			}
		}

		inputName := field.Name

		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			inputName = strings.Split(jsonTag, ",")[0]
		}

		minTag := field.Tag.Get("min")
		maxTag := field.Tag.Get("max")
		requiredTag := field.Tag.Get("required")

		fields = append(fields, DynamicConfigField{
			Name:     inputName,
			Label:    label,
			Type:     inputType,
			Help:     helpText,
			Required: requiredTag == "true",
			Min:      minTag,
			Max:      maxTag,
			Options:  options,
			Default:  defaultValue,
		})
	}
	return fields
}

func ConvertToDynamicConfig(target DynamicConfig, config any) (DynamicConfig, error) {
	typedConfig := reflect.New(reflect.TypeOf(target)).Interface()
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
	return typedConfig.(DynamicConfig), nil
}
