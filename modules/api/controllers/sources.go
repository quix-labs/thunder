package controllers

import (
	"fmt"
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/modules/http_server"
	"github.com/quix-labs/thunder/modules/http_server/helpers"
	"net/http"
)

func SourceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /sources", listSources)
	mux.HandleFunc("POST /sources", createSource)
	mux.HandleFunc("PUT /sources/{id}", updateSource)
	mux.HandleFunc("DELETE /sources/{id}", deleteSource)
	mux.HandleFunc("GET /sources/{id}/stats", getSourceStats)
}

type SourceApiDetails struct {
	ID      string `json:"id"`
	Driver  string `json:"driver"`
	Config  any    `json:"config"`
	Excerpt string `json:"excerpt"`
}

func listSources(w http.ResponseWriter, r *http.Request) {
	sources := thunder.Sources.All()
	var res []SourceApiDetails

	for _, source := range sources {
		serializeSource := helpers.Must(thunder.SerializeSource(&source))
		res = append(res, SourceApiDetails{
			ID:      serializeSource.ID,
			Driver:  serializeSource.Driver,
			Config:  serializeSource.Config,
			Excerpt: source.Config.Excerpt(),
		})
	}
	helpers.WriteJsonResponse(w, http.StatusOK, res)
}

func createSource(w http.ResponseWriter, r *http.Request) {
	var p thunder.JsonSource
	helpers.CheckErr(http_server.DecodeJSONBody(w, r, &p))

	helpers.NextCheckStatus(http.StatusBadRequest)
	source := helpers.Must(thunder.UnserializeSource(&p))

	helpers.CheckErr(thunder.Sources.Register("", *source))

	helpers.CheckErr(thunder.SaveConfig())

	helpers.WriteJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, "Source created"})
}

func updateSource(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var s thunder.JsonSource
	helpers.CheckErr(http_server.DecodeJSONBody(w, r, &s))

	newSource := helpers.Must(thunder.UnserializeSource(&s))

	helpers.CheckErr(thunder.Sources.Update(id, *newSource))
	helpers.CheckErr(thunder.SaveConfig())

	helpers.WriteJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, fmt.Sprintf("Source %s updated", id)})
}

func deleteSource(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	helpers.CheckErr(thunder.Sources.Delete(id))
	helpers.CheckErr(thunder.SaveConfig())
	fmt.Println(thunder.Sources.All())
	helpers.WriteJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, fmt.Sprintf(`Source %s deleted!`, id)})
}

func getSourceStats(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	helpers.NextCheckStatus(http.StatusBadRequest)
	source := helpers.Must(thunder.Sources.Get(id))

	helpers.NextCheckStatus(http.StatusUnprocessableEntity)
	stats := helpers.Must(source.Driver.Stats())

	helpers.WriteJsonResponse(w, http.StatusOK, stats)
}
