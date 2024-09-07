package http_server

import (
	"errors"
	"fmt"
	"github.com/quix-labs/thunder"
	"net/http"
)

func init() {
	thunder.RegisterModule(&Module{})
}

type Module struct{}

func (m *Module) ThunderModule() thunder.ModuleInfo {
	return thunder.ModuleInfo{
		ID:  "thunder.http_server",
		New: func() thunder.Module { return new(Module) },
	}
}

func (m *Module) Start() error {
	if !IsHttpServerEnabled() {
		return nil
	}

	srv := &http.Server{
		Addr:    GetHttpServerAddr(),
		Handler: router(),
	}

	errChan := make(chan error, 1)
	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
		errChan <- nil
	}()

	fmt.Println("http server listening on", srv.Addr)
	return <-errChan
}
func (m *Module) Stop() error {
	return nil
}

var (
	_ thunder.Module = (*Module)(nil) // Interface implementation
)

// Allow external configuration

var (
	httpServerEnabled = true
	httpServerAddr    = ":3000"
)

func SetHttpServerEnabled(enabled bool) {
	httpServerEnabled = enabled
}
func IsHttpServerEnabled() bool {
	return httpServerEnabled
}

func SetHttpServerAddr(addr string) {
	httpServerAddr = addr
}
func GetHttpServerAddr() string {
	return httpServerAddr
}

// Internal extending

func router() http.Handler {
	mux := http.NewServeMux()

	for _, handler := range GetHandlers() {
		handler.HandleRoutes(mux)
	}

	return mux
}
