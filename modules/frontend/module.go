package frontend

import (
	"errors"
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/modules/http_server"
	"net/http"
)

const ModuleID = "thunder.frontend"

func init() {
	thunder.Modules.Register(ModuleID, &Module{})
	http_server.RegisterHandler(&Module{})
}

type Module struct{}

func (m *Module) RequiredModules() []string {
	return []string{
		"thunder.http_server",
		"thunder.api",
	}
}

func (m *Module) New() thunder.Module {
	return new(Module)
}

func (m *Module) HandleRoutes(mux *http.ServeMux) {
	HandleFrontend(mux)
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
