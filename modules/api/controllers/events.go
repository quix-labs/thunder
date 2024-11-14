package controllers

import (
	"net/http"
)

func EventsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /events", events)
}

func events(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers to allow all origins. You may want to restrict this to specific origins in a production environment.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	//TODO STREAM EVENTS
	//stopListening := thunder.GetEventBroadcaster().OnAll(func(event string, data any) {
	//	payload, err := json.Marshal(data)
	//	if err != nil {
	//		return // TODO ERROR
	//	}
	//	_, err = fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, payload)
	//	if err != nil {
	//		return // TODO ERROR
	//	}
	//	w.(http.Flusher).Flush()
	//})

	<-r.Context().Done()
	//stopListening()

}
