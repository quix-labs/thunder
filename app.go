package thunder

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
)

func Start() error {
	err := LoadConfig()
	if err != nil {
		return err
	}

	// Run moduleInfo in parallel
	moduleErrChan := make(chan error)
	for _, moduleInfo := range GetModules() {
		// Check dependencies
		for _, requiredModule := range moduleInfo.RequiredModules {
			if _, err := GetModule(requiredModule); err != nil {
				return fmt.Errorf(`module "%s" required module "%s" to be present`, moduleInfo.ID, requiredModule)
			}
		}
		module := moduleInfo.New()
		go func() {
			err := module.Start()
			if err != nil {
				moduleErrChan <- err
			}
		}()
	}

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
