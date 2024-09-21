package thunder

import (
	"github.com/rs/zerolog"
	"os"
)

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

func GetLoggerForModule(m Module) *zerolog.Logger {
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel).With().Str("module", m.ThunderModule().ID).Stack().Timestamp().Logger()
	return &logger
}
