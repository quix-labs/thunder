package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/modules/http_server"
	"log"
	"net/http"
	"reflect"
	"strings"
)

func SourceDriverRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /go-api/source-drivers", listSourceDrivers)
	mux.HandleFunc("POST /go-api/source-drivers/test", testSourceDriver)
}

func testSourceDriver(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var p struct {
		Driver string
		Config map[string]any
	}

	err := http_server.DecodeJSONBody(w, r, &p)
	if err != nil {
		var mr *http_server.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	driver, err := thunder.GetSourceDriver(p.Driver)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"%s"}`, err)))
		return
	}

	configInstance, err := thunder.ConvertSourceDriverConfig(&driver, p.Config)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"%s"}`, err)))
		return
	}

	// TRY TEST
	message, err := driver.New().TestConfig(configInstance)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		bytes, _ := json.Marshal(struct {
			Success bool   `json:"success"`
			Error   string `json:"error"`
			Message string `json:"message"`
		}{Success: false, Error: err.Error(), Message: message})
		w.Write(bytes)
		return
	}
	w.Write([]byte(fmt.Sprintf(`{"success":true,"message":"%s"}`, message)))
}

type FieldDetails struct {
	Name     string  `json:"name"`
	Label    string  `json:"label"`
	Type     string  `json:"type"`
	Required bool    `json:"required"`
	Help     *string `json:"help,omitempty"`
	Default  *string `json:"default,omitempty"`
}

type DriverDetails struct {
	Config *thunder.SourceDriverInfo `json:"config"`
	Fields []FieldDetails            `json:"fields"`
}

func listSourceDrivers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	registeredDrivers := thunder.GetSourceDrivers()
	res := make(map[string]*DriverDetails, len(registeredDrivers))

	for _, driver := range registeredDrivers {
		res[driver.ID] = &DriverDetails{
			Config: &driver,
			Fields: parseSourceDriverFields(&driver),
		}
	}

	jsonData, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"%s"}`, err)))
		return
	}

	w.Write(jsonData)
}

func parseSourceDriverFields(driver *thunder.SourceDriverInfo) []FieldDetails {

	var fields []FieldDetails
	configValue := reflect.ValueOf(driver.Config)
	configType := reflect.TypeOf(driver.Config)

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
			Default:  defaultValue,
		})
	}
	return fields
}
