package controllers

import (
	"errors"
	"fmt"
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/modules/http_server"
	"log"
	"net/http"
	"strconv"
)

func TargetRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /go-api/targets", listTargets)
	mux.HandleFunc("POST /go-api/targets", createTarget)
	mux.HandleFunc("DELETE /go-api/targets/{id}", deleteTarget)
}

func deleteTarget(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, errors.New("IS is not an integer"), "")
		return
	}

	config := thunder.GetConfig()
	config.Targets = append(config.Targets[:id], config.Targets[id+1:]...)
	thunder.SetConfig(config)

	err = thunder.SaveConfig()
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	writeJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, fmt.Sprintf(`{"success":true,"message":"Source %d deleted!"}`, id)})
}

func listTargets(w http.ResponseWriter, r *http.Request) {
	writeJsonResponse(w, http.StatusOK, thunder.GetConfig().Targets)
}

func createTarget(w http.ResponseWriter, r *http.Request) {
	var p struct {
		Driver string `json:"driver"`
		Config any    `json:"config"`
	}
	err := http_server.DecodeJSONBody(w, r, &p)
	if err != nil {
		var mr *http_server.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	driver, err := thunder.GetTargetDriver(p.Driver)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err, "")
		return
	}

	configInstance, err := thunder.ConvertTargetDriverConfig(&driver, p.Config)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	config := thunder.GetConfig()
	config.Targets = append(config.Targets, thunder.Target{Driver: p.Driver, Config: configInstance})
	thunder.SetConfig(config)

	err = thunder.SaveConfig()
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}
	writeJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, "Target created"})
}
