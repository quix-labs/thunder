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

	// Run http server
	httpServerErrChan := make(chan error)
	go func() {
		err := StartHttpServer()
		if err != nil {
			httpServerErrChan <- err
		}
	}()

	// Run module in parallel
	moduleErrChan := make(chan error)
	for _, module := range GetModules() {
		module := module.New()
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
	case err := <-httpServerErrChan:
		return err
	case err := <-moduleErrChan:
		return err
	}

	return nil
}
