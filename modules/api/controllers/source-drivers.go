package controllers

import (
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/modules/http_server"
	"github.com/quix-labs/thunder/modules/http_server/helpers"
	"github.com/quix-labs/thunder/utils"
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
	helpers.CheckErr(http_server.DecodeJSONBody(w, r, &p))

	helpers.NextCheckStatus(http.StatusBadRequest)
	driver := helpers.Must(thunder.SourceDrivers.Get(p.Driver))

	config := driver.Config().Config
	configInstance := helpers.Must(utils.ConvertToDynamicConfig(config, p.Config))

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

type DriverDetails struct {
	Config *thunder.SourceDriverConfig `json:"config"`
	Fields utils.DynamicConfigFields   `json:"fields"`
}

func listSourceDrivers(w http.ResponseWriter, r *http.Request) {
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
