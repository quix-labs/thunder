package controllers

import (
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/modules/http_server"
	"github.com/quix-labs/thunder/modules/http_server/helpers"
	"github.com/quix-labs/thunder/utils"
	"net/http"
)

func TargetDriverRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /target-drivers", listTargetDrivers)
	mux.HandleFunc("POST /target-drivers/test", testTargetDriver)
}

func testTargetDriver(w http.ResponseWriter, r *http.Request) {
	var p struct {
		Driver string
		Config map[string]any
	}
	helpers.CheckErr(http_server.DecodeJSONBody(w, r, &p))

	helpers.NextCheckStatus(http.StatusBadRequest)
	driver := helpers.Must(thunder.TargetDrivers.Get(p.Driver))

	driverConfig := driver.Config()
	configInstance := helpers.Must(utils.ConvertToDynamicConfig(driverConfig.Config, p.Config))

	// TRY TEST
	helpers.NextCheckStatus(http.StatusUnprocessableEntity)
	driverInstance := helpers.Must(driver.New(configInstance))

	helpers.NextCheckStatus(http.StatusBadRequest)
	message := helpers.Must(driverInstance.TestConfig())

	helpers.WriteJsonResponse(w, http.StatusOK, struct {
		Success bool
		Message string `json:"message"`
	}{true, message})
}

type TargetDriverDetails struct {
	Config *thunder.TargetDriverConfig `json:"config"`
	Fields utils.DynamicConfigFields   `json:"fields"`
}

func listTargetDrivers(w http.ResponseWriter, r *http.Request) {
	registeredDrivers := thunder.TargetDrivers.All()
	res := make(map[string]*TargetDriverDetails, len(registeredDrivers))
	for driverID, driver := range registeredDrivers {
		driverConfig := driver.Config()

		res[driverID] = &TargetDriverDetails{
			Config: &driverConfig,
			Fields: utils.ParseDynamicConfigFields(&driverConfig.Config),
		}
	}
	helpers.WriteJsonResponse(w, http.StatusOK, res)
}
