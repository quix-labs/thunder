package thunder

type ModuleInfo struct {
	ID              string
	New             func() Module
	RequiredModules []string
}

type Module interface {
	ThunderModule() ModuleInfo

	Start() error
	Stop() error
}
