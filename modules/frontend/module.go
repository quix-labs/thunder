package frontend

import "github.com/quix-labs/thunder"

func init() {
	thunder.RegisterModule(&Module{})
}

type Module struct{}

func (m *Module) ThunderModule() thunder.ModuleInfo {
	return thunder.ModuleInfo{
		ID:  "frontend",
		New: func() thunder.Module { return new(Module) },
	}
}

func (m *Module) Start() error {
	return Start()
}

func (m *Module) Stop() error {
	return nil
}

var (
	_ thunder.Module = (*Module)(nil) // Interface implementation
)
