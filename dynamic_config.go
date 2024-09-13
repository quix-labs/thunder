package thunder

import (
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
	Default  *string `json:"default,omitempty"`
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

		label := field.Tag.Get("label")
		inputType := field.Tag.Get("type")
		helpTag := field.Tag.Get("help")
		defaultTag := field.Tag.Get("default")
		requiredTag := field.Tag.Get("required")
		minTag := field.Tag.Get("min")

		if label == "" {
			label = field.Name
		}
		if inputType == "" {
			inputType = "text"
		}
		var helpText *string = nil
		if helpTag != "" {
			helpText = &helpTag
		}

		var defaultValue *string = nil
		if defaultTag != "" {
			defaultValue = &defaultTag
		}

		inputName := field.Name

		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			inputName = strings.Split(jsonTag, ",")[0]
		}

		fields = append(fields, DynamicConfigField{
			Name:     inputName,
			Label:    label,
			Type:     inputType,
			Help:     helpText,
			Required: requiredTag == "true",
			Min:      minTag,
			Default:  defaultValue,
		})
	}
	return fields
}

func ConvertToDynamicConfig(target *DynamicConfig, config any) (DynamicConfig, error) {
	typedConfig := reflect.New(reflect.TypeOf(*target)).Interface()
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
