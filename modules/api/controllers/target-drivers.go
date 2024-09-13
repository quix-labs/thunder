package controllers

import (
	"errors"
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/modules/http_server"
	"log"
	"net/http"
)

func TargetDriverRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /go-api/target-drivers", listTargetDrivers)
	mux.HandleFunc("POST /go-api/target-drivers/test", testTargetDriver)
}

func testTargetDriver(w http.ResponseWriter, r *http.Request) {
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
	driver, err := thunder.GetTargetDriver(p.Driver)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err, "")
		return
	}

	configInstance, err := thunder.ConvertToDynamicConfig(&driver.Config, p.Config)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	// TRY TEST
	driverInstance, err := driver.New(configInstance)
	if err != nil {
		writeJsonError(w, http.StatusUnprocessableEntity, err, "")
		return
	}

	message, err := driverInstance.TestConfig()
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err, message)
		return
	}
	writeJsonResponse(w, http.StatusOK, struct {
		Success bool
		Message string `json:"message"`
	}{true, message})
}

type TargetDriverDetails struct {
	Config *thunder.TargetDriverInfo   `json:"config"`
	Fields thunder.DynamicConfigFields `json:"fields"`
}

func listTargetDrivers(w http.ResponseWriter, r *http.Request) {
	registeredDrivers := thunder.GetTargetDrivers()
	res := make(map[string]*TargetDriverDetails, len(registeredDrivers))
	for _, driver := range registeredDrivers {
		res[driver.ID] = &TargetDriverDetails{
			Config: &driver,
			Fields: thunder.ParseDynamicConfigFields(&driver.Config),
		}
	}
	writeJsonResponse(w, http.StatusOK, res)
}
