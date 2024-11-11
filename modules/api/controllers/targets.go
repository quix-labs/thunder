package controllers

import (
	"fmt"
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/modules/http_server"
	"github.com/quix-labs/thunder/modules/http_server/helpers"
	"net/http"
)

func TargetRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /targets", listTargets)
	mux.HandleFunc("POST /targets", createTarget)
	mux.HandleFunc("PUT /targets/{id}", updateTarget)
	mux.HandleFunc("DELETE /targets/{id}", deleteTarget)
}

type TargetApiDetails struct {
	ID      string `json:"id"`
	Driver  string `json:"driver"`
	Config  any    `json:"config"`
	Excerpt string `json:"excerpt"`
}

func listTargets(w http.ResponseWriter, r *http.Request) {
	var res []TargetApiDetails

	for _, target := range thunder.Targets.All() {
		serializeTarget := helpers.Must(thunder.SerializeTarget(&target))
		res = append(res, TargetApiDetails{
			ID:      serializeTarget.ID,
			Driver:  serializeTarget.Driver,
			Config:  serializeTarget.Config,
			Excerpt: target.Config.Excerpt(),
		})
	}
	helpers.WriteJsonResponse(w, http.StatusOK, res)
}

func createTarget(w http.ResponseWriter, r *http.Request) {
	var p thunder.JsonTarget
	helpers.CheckErr(http_server.DecodeJSONBody(w, r, &p))

	helpers.NextCheckStatus(http.StatusBadRequest)
	target := helpers.Must(thunder.UnserializeTarget(&p))

	helpers.CheckErr(thunder.Targets.Register("", *target))
	helpers.CheckErr(thunder.SaveConfig())

	helpers.WriteJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, "Target created"})
}

func updateTarget(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var s thunder.JsonTarget
	helpers.CheckErr(http_server.DecodeJSONBody(w, r, &s))

	newTarget := helpers.Must(thunder.UnserializeTarget(&s))

	helpers.CheckErr(thunder.Targets.Update(id, *newTarget))
	helpers.CheckErr(thunder.SaveConfig())

	helpers.WriteJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, fmt.Sprintf("Target %s updated", id)})
}

func deleteTarget(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	helpers.CheckErr(thunder.Targets.Delete(id))
	helpers.CheckErr(thunder.SaveConfig())
	helpers.WriteJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, fmt.Sprintf(`Target %s deleted!`, id)})
}
