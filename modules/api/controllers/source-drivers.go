package controllers

import (
	"errors"
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/modules/http_server"
	"github.com/quix-labs/thunder/modules/http_server/helpers"
	"github.com/quix-labs/thunder/utils"
	"log"
	"net/http"
)

func SourceDriverRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /source-drivers", listSourceDrivers)
	mux.HandleFunc("POST /source-drivers/test", testSourceDriver)
}

func testSourceDriver(w http.ResponseWriter, r *http.Request) {
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
	driver, err := thunder.SourceDrivers.Get(p.Driver)
	if err != nil {
		helpers.WriteJsonError(w, http.StatusBadRequest, err, "")
		return
	}

	config := driver.Config().Config
	configInstance, err := utils.ConvertToDynamicConfig(config, p.Config)
	if err != nil {
		helpers.WriteJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	// TRY TEST
	driverInstance, err := driver.New(configInstance)
	if err != nil {
		helpers.WriteJsonError(w, http.StatusUnprocessableEntity, err, "")
		return
	}

	message, err := driverInstance.TestConfig()
	if err != nil {
		helpers.WriteJsonError(w, http.StatusBadRequest, err, message)
		return
	}
	helpers.WriteJsonResponse(w, http.StatusOK, struct {
		Success bool
		Message string `json:"message"`
	}{true, message})
}

type DriverDetails struct {
	Config *thunder.SourceDriverConfig `json:"config"`
	Fields utils.DynamicConfigFields   `json:"fields"`
}

func listSourceDrivers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	registeredDrivers := thunder.SourceDrivers.All()
	res := make(map[string]*DriverDetails, len(registeredDrivers))

	for driverID, driver := range registeredDrivers {
		driverConfig := driver.Config()

		res[driverID] = &DriverDetails{
			Config: &driverConfig,
			Fields: utils.ParseDynamicConfigFields(&driverConfig.Config),
		}
	}
	helpers.WriteJsonResponse(w, http.StatusOK, res)
}
