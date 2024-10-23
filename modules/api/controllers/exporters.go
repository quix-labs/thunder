package controllers

import (
	"github.com/quix-labs/thunder"
	"net/http"
)

func ExporterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /go-api/exporters", listExporters)
}

func listExporters(w http.ResponseWriter, r *http.Request) {
	exporters := thunder.Exporters.All()
	var response = make(map[string]string, len(exporters))
	for exporterID, exporter := range exporters {
		response[exporterID] = exporter.Name()
	}
	writeJsonResponse(w, http.StatusOK, response)
}
