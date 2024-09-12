package thunder

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
)

type ProcessorState string

// TODO DECOUPER LES LOGIQUE DE STATE DANS CHAQUE FICHIERS RESPECTIFS
var (
	ProcessorIndexing = ProcessorState("indexing")
	ProcessorActive   = ProcessorState("active")
	ProcessorInactive = ProcessorState("inactive")
)

type LoadedProcessors map[int]struct { // As [index]struct
	State        ProcessorState `json:"state"`
	SourceDriver SourceDriver   `json:"-"`
}

type AppState struct {
	LoadedProcessor LoadedProcessors `json:"loaded_processors"`
	LoadedModules   []Module         `json:"-"`
}

var state AppState

func Start() error {
	err := LoadConfig()
	if err != nil {
		return err
	}

	// Load Modules
	modules := GetModules()
	state.LoadedModules = make([]Module, len(modules))
	moduleErrChan := make(chan error)
	for index, moduleInfo := range modules {
		// Check dependencies
		for _, requiredModule := range moduleInfo.RequiredModules {
			if _, err := GetModule(requiredModule); err != nil {
				return fmt.Errorf(`module "%s" depends on missing module "%s"`, moduleInfo.ID, requiredModule)
			}
		}

		state.LoadedModules[index] = moduleInfo.New()
		go func() {
			err := state.LoadedModules[index].Start()
			if err != nil {
				moduleErrChan <- err
			}
		}()
	}

	//processors := GetConfig().Processors

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-ctx.Done():
		fmt.Println("TODO GRACEFULL SHUTDOWN")
	case err := <-moduleErrChan:
		return err
	}

	return nil
}
