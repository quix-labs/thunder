package controllers

import (
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/modules/http_server/helpers"
	"net/http"
)

func ExporterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /exporters", listExporters)
}

func listExporters(w http.ResponseWriter, r *http.Request) {
	exporters := thunder.Exporters.All()
	var response = make(map[string]string, len(exporters))
	for exporterID, exporter := range exporters {
		response[exporterID] = exporter.Name()
	}
	helpers.WriteJsonResponse(w, http.StatusOK, response)
}
