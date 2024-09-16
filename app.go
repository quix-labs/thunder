package thunder

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
)

type App struct {
	loadedModules []Module

	// Used by module to handle app events
	broadcaster any
}

var app = new(App)

func GetApp() *App {
	return app
}

func Start() error {
	app := GetApp()

	err := LoadConfig()
	if err != nil {
		return err
	}

	// Load Modules
	modules := GetModules()
	app.loadedModules = make([]Module, len(modules))
	moduleErrChan := make(chan error)
	for index, moduleInfo := range modules {
		// Check dependencies
		for _, requiredModule := range moduleInfo.RequiredModules {
			if _, err := GetModule(requiredModule); err != nil {
				return fmt.Errorf(`module "%s" depends on missing module "%s"`, moduleInfo.ID, requiredModule)
			}
		}

		app.loadedModules[index] = moduleInfo.New()
		go func() {
			err := app.loadedModules[index].Start()
			if err != nil {
				moduleErrChan <- err
			}
		}()
	}
	//
	//var processorsWg sync.WaitGroup
	//var processorErrChan = make(chan error)
	//
	//for _, p := range GetProcessors() {
	//	//if !p.Enabled {
	//	//	continue
	//	//}
	//	processorsWg.Add(1)
	//	go func() {
	//		defer func() {
	//			processorsWg.Done()
	//			fmt.Println("processor stopped")
	//		}()
	//		err := p.Start()
	//		if err != nil && !errors.Is(err, context.Canceled) {
	//			processorErrChan <- err
	//		}
	//	}()
	//}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-ctx.Done():
		fmt.Println("TODO GRACEFULL SHUTDOWN")
	//case err := <-processorErrChan:
	//	return err
	case err := <-moduleErrChan:
		return err
	}

	// STOP ALL PROCESSORS
	for _, p := range GetProcessors() {
		err := p.Stop()
		if err != nil {
			panic(err)
		}
	}

	//processorsWg.Wait()

	return nil
}
