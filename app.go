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

	// Load Processors
	//if err = LoadAllProcessors(); err != nil {
	//	return err
	//}

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

func (app *App) Reload() error {
	panic("not implemented")
}

func (app *App) IndexProcessor() error {
	panic("not implemented")
}

func (app *App) AddProcessor() error {
	panic("not implemented")
}

func (app *App) GetProcessor(id int) error {
	panic("not implemented")
}

func (app *App) GetProcessors(id int) error {
	panic("not implemented")
}

func (app *App) LoadConfig() error {
	panic("not implemented")
}

func (app *App) SaveConfig() error {
	panic("not implemented")
}
