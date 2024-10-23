package thunder

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
)

type App struct {
	loadedModules map[string]Module

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
	modules := Modules.All()
	app.loadedModules = make(map[string]Module, len(modules))

	moduleErrChan := make(chan error)
	for moduleID, module := range modules {
		// Check dependencies
		for _, requiredModule := range module.RequiredModules() {
			if _, err := Modules.Get(requiredModule); err != nil {
				return fmt.Errorf(`module "%s" depends on missing module "%s"`, moduleID, requiredModule)
			}
		}
		moduleInstance := module.New()
		go func() {
			err := moduleInstance.Start()
			if err != nil {
				moduleErrChan <- err
			}
		}()
		app.loadedModules[moduleID] = moduleInstance
	}

	//var processorsWg sync.WaitGroup
	//var processorErrChan = make(chan error)
	//
	//for _, p := range GetProcessors() {
	//	//if !p.Enabled {
	//	//	continue
	//	//}
	//	processorsWg.Add(1)
	//	go func() {
	//		defer processorsWg.Done()
	//
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
	//case err := <-processorErrChan:
	//	return err
	case err := <-moduleErrChan:
		return err
	}

	// STOP ALL PROCESSORS
	for _, p := range Processors.All() {
		err := p.Stop()
		if err != nil {
			panic(err)
		}
	}

	// STOP ALL MODULES
	for _, p := range app.loadedModules {
		err := p.Stop()
		if err != nil {
			panic(err)
		}
	}

	//processorsWg.Wait()

	return nil
}
