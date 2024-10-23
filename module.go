package thunder

import (
	"errors"
	"github.com/quix-labs/thunder/utils"
	"github.com/rs/zerolog"
	"os"
)

type Module interface {
	RequiredModules() []string

	New() Module

	Start() error
	Stop() error
}

func GetLoggerForModule(ModuleID string) *zerolog.Logger {
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel).With().Str("module", ModuleID).Stack().Timestamp().Logger()
	return &logger
}

// Modules is a registry that allows external library to register their own modules
var Modules = utils.NewRegistry[Module]("module").ValidateUsing(func(ID string, module Module) error {
	if module.New == nil {
		return errors.New(`module doesn't implements "New" method`)
	}
	if val := module.New(); val == nil {
		return errors.New(`module "New" method must return a non-nil instance`)
	}
	return nil
})
