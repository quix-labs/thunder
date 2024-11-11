package http_server

import (
	"errors"
	"github.com/quix-labs/thunder"
	"net/http"
)

const ModuleID = "thunder.http_server"

func init() {
	thunder.Modules.Register(ModuleID, &Module{})
}

type Module struct{}

func (m *Module) RequiredModules() []string {
	return []string{}
}

func (m *Module) New() thunder.Module {
	return new(Module)
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

	thunder.GetLoggerForModule(ModuleID).Info().Msgf("start listening on %s", srv.Addr)

	return <-errChan
}
func (m *Module) Stop() error {
	thunder.GetLoggerForModule(ModuleID).Info().Msg("stop listening")
	return nil
}

var (
	_ thunder.Module = (*Module)(nil)
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
