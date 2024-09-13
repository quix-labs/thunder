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
	mux.HandleFunc("PUT /go-api/targets/{id}", updateTarget)
	mux.HandleFunc("DELETE /go-api/targets/{id}", deleteTarget)
}

type TargetApiDetails struct {
	Driver  string `json:"driver"`
	Config  any    `json:"config"`
	Excerpt string `json:"excerpt"`
}

func listTargets(w http.ResponseWriter, r *http.Request) {
	targets := thunder.GetConfig().Targets
	var res = make(map[int]TargetApiDetails, len(targets))

	for key, target := range targets {
		res[key] = TargetApiDetails{
			Driver:  target.Driver,
			Config:  target.Config,
			Excerpt: target.Config.Excerpt(),
		}
	}
	writeJsonResponse(w, http.StatusOK, res)
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

	configInstance, err := thunder.ConvertToDynamicConfig(&driver.Config, p.Config)
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

func updateTarget(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, errors.New("ID is not an integer"), "")
		return
	}

	var p struct {
		Driver string `json:"driver"`
		Config any    `json:"config"`
	}
	err = http_server.DecodeJSONBody(w, r, &p)
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

	configInstance, err := thunder.ConvertToDynamicConfig(&driver.Config, p.Config)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	if (id + 1) > len(thunder.GetConfig().Targets) {
		writeJsonError(w, http.StatusBadRequest, errors.New("invalid target"), "")
		return
	}

	config := thunder.GetConfig()
	config.Targets[id] = thunder.Target{Driver: p.Driver, Config: configInstance}
	thunder.SetConfig(config)

	err = thunder.SaveConfig()
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}
	writeJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, fmt.Sprintf("Target %d updated", id)})
}

func deleteTarget(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, errors.New("ID is not an integer"), "")
		return
	}

	if (id + 1) > len(thunder.GetConfig().Targets) {
		writeJsonError(w, http.StatusBadRequest, errors.New("invalid target"), "")
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
	}{true, fmt.Sprintf(`Target %d deleted!`, id)})
}
