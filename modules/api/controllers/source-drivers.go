package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/modules/http_server"
	"log"
	"net/http"
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
	driverInstance, err := driver.New(configInstance)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		bytes, _ := json.Marshal(struct {
			Success bool   `json:"success"`
			Error   string `json:"error"`
		}{Success: false, Error: err.Error()})
		w.Write(bytes)
		return
	}

	message, err := driverInstance.TestConfig()
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
			Fields: parseConfigFields(driver.Config),
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
