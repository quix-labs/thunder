package api

import (
	"errors"
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/modules/api/controllers"
	"github.com/quix-labs/thunder/modules/http_server"
	"net/http"
)

func init() {
	thunder.RegisterModule(&Module{})
	http_server.RegisterHandler(&Module{})
}

type Module struct{}

func (m *Module) ThunderModule() thunder.ModuleInfo {
	return thunder.ModuleInfo{
		ID:              "thunder.api",
		New:             func() thunder.Module { return new(Module) },
		RequiredModules: []string{"thunder.http_server"},
	}
}

func (m *Module) HandleRoutes(mux *http.ServeMux) {
	controllers.SourceRoutes(mux)
	controllers.SourceDriverRoutes(mux)
	controllers.ProcessorRoutes(mux)
	controllers.TargetDriverRoutes(mux)
	controllers.TargetRoutes(mux)
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
