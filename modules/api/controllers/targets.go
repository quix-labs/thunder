package controllers

import (
	"errors"
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
		serializeTarget, err := thunder.SerializeTarget(&target)
		if err != nil {
			helpers.WriteJsonError(w, http.StatusInternalServerError, err, "")
			return
		}

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
	err := http_server.DecodeJSONBody(w, r, &p)
	if err != nil {
		var mr *http_server.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			thunder.GetLoggerForModule("thunder.api").Error().Msg(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	target, err := thunder.UnserializeTarget(&p)
	if err != nil {
		helpers.WriteJsonError(w, http.StatusBadRequest, err, "")
		return
	}

	err = thunder.Targets.Register("", *target)
	if err != nil {
		helpers.WriteJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	err = thunder.SaveConfig()
	if err != nil {
		helpers.WriteJsonError(w, http.StatusInternalServerError, err, "")
		return
	}
	helpers.WriteJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, "Target created"})
}

func updateTarget(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var s thunder.JsonTarget
	err := http_server.DecodeJSONBody(w, r, &s)
	if err != nil {
		var mr *http_server.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			thunder.GetLoggerForModule("thunder.api").Error().Msg(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	newTarget, err := thunder.UnserializeTarget(&s)
	if err != nil {
		helpers.WriteJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	if err := thunder.Targets.Update(id, *newTarget); err != nil {
		helpers.WriteJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	if err = thunder.SaveConfig(); err != nil {
		helpers.WriteJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	helpers.WriteJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, fmt.Sprintf("Target %d updated", id)})
}

func deleteTarget(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := thunder.Targets.Delete(id)
	if err != nil {
		helpers.WriteJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	err = thunder.SaveConfig()
	if err != nil {
		helpers.WriteJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	helpers.WriteJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, fmt.Sprintf(`Target %d deleted!`, id)})
}
