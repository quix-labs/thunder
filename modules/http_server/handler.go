package http_server

import (
	"net/http"
	"sync"
)

type Handler interface {
	HandleRoutes(mux *http.ServeMux)
}

// REGISTERING HANDLERS
var (
	handlers   []Handler
	handlersMu sync.RWMutex
)

func RegisterHandler(handler Handler) {
	handlersMu.Lock()
	defer handlersMu.Unlock()
	handlers = append(handlers, handler)
}

func GetHandlers() []Handler {
	handlersMu.RLock()
	defer handlersMu.RUnlock()
	return handlers
}
