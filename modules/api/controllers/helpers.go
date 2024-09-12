package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type FieldDetails struct {
	Name     string  `json:"name"`
	Label    string  `json:"label"`
	Type     string  `json:"type"`
	Required bool    `json:"required"`
	Help     *string `json:"help,omitempty"`
	Min      string  `json:"min,omitempty"`
	Default  *string `json:"default,omitempty"`
}

func parseConfigFields(config any) []FieldDetails {

	var fields []FieldDetails
	configValue := reflect.ValueOf(config)
	configType := reflect.TypeOf(config)

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

		fields = append(fields, FieldDetails{
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

func writeJsonResponse(w http.ResponseWriter, statusCode int, response interface{}) {

	marshal, err := json.Marshal(response)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(marshal)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
	}
}

func writeJsonError(w http.ResponseWriter, statusCode int, error error, message string) {
	var payload = struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
		Message string `json:"message"`
	}{
		Success: false,
		Error:   error.Error(),
		Message: message,
	}
	marshal, err := json.Marshal(&payload)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(marshal)
	if err != nil {
		log.Println(err)
	}
}
