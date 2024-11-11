package api

import (
	"errors"
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/modules/api/controllers"
	"github.com/quix-labs/thunder/modules/http_server"
	"github.com/quix-labs/thunder/modules/http_server/helpers"
	"net/http"
)

const ModuleID = "thunder.api"

func init() {
	thunder.Modules.Register(ModuleID, &Module{})
	http_server.RegisterHandler(&Module{})
}

type Module struct{}

func (m *Module) RequiredModules() []string {
	return []string{"thunder.http_server"}
}

func (m *Module) New() thunder.Module {
	return new(Module)
}

func (m *Module) HandleRoutes(mux *http.ServeMux) {
	innerMux := http.NewServeMux()

	controllers.SourceRoutes(innerMux)
	controllers.SourceDriverRoutes(innerMux)
	controllers.ProcessorRoutes(innerMux)
	controllers.TargetDriverRoutes(innerMux)
	controllers.TargetRoutes(innerMux)
	controllers.EventsRoutes(innerMux)
	controllers.ExporterRoutes(innerMux)

	mux.Handle("/go-api/", helpers.ErrorMiddleware(http.StripPrefix("/go-api", innerMux), ModuleID))
}

func (m *Module) Start() error {
	if !http_server.IsHttpServerEnabled() {
		return errors.New("http_server need to be enabled to run frontend")
	}
	return nil
}

func (m *Module) Stop() error {
	return nil
}

var (
	_ thunder.Module      = (*Module)(nil) // Interface implementation
	_ http_server.Handler = (*Module)(nil) // Interface implementation
)
