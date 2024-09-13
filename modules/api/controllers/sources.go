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

func SourceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /go-api/sources", listSources)
	mux.HandleFunc("POST /go-api/sources", createSource)
	mux.HandleFunc("PUT /go-api/sources/{id}", updateSource)
	mux.HandleFunc("DELETE /go-api/sources/{id}", deleteSource)
	mux.HandleFunc("GET /go-api/sources/{id}/stats", getSourceStats)
}

type SourceApiDetails struct {
	Driver  string `json:"driver"`
	Config  any    `json:"config"`
	Excerpt string `json:"excerpt"`
}

func listSources(w http.ResponseWriter, r *http.Request) {
	sources := thunder.GetConfig().Sources
	var res = make(map[int]SourceApiDetails, len(sources))

	for key, source := range sources {
		res[key] = SourceApiDetails{
			Driver:  source.Driver,
			Config:  source.Config,
			Excerpt: source.Config.Excerpt(),
		}
	}
	writeJsonResponse(w, http.StatusOK, res)
}

func createSource(w http.ResponseWriter, r *http.Request) {
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

	driver, err := thunder.GetSourceDriver(p.Driver)
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
	config.Sources = append(config.Sources, thunder.Source{Driver: p.Driver, Config: configInstance})
	thunder.SetConfig(config)

	err = thunder.SaveConfig()
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}
	writeJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, "Source created"})
}

func updateSource(w http.ResponseWriter, r *http.Request) {
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

	driver, err := thunder.GetSourceDriver(p.Driver)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err, "")
		return
	}

	configInstance, err := thunder.ConvertToDynamicConfig(&driver.Config, p.Config)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	if (id + 1) > len(thunder.GetConfig().Sources) {
		writeJsonError(w, http.StatusBadRequest, errors.New("invalid source"), "")
		return
	}

	config := thunder.GetConfig()
	config.Sources[id] = thunder.Source{Driver: p.Driver, Config: configInstance}
	thunder.SetConfig(config)

	err = thunder.SaveConfig()
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}
	writeJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, fmt.Sprintf("Source %d updated", id)})
}

func deleteSource(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, errors.New("ID is not an integer"), "")
		return
	}

	if (id + 1) > len(thunder.GetConfig().Sources) {
		writeJsonError(w, http.StatusBadRequest, errors.New("invalid source"), "")
		return
	}

	config := thunder.GetConfig()
	config.Sources = append(config.Sources[:id], config.Sources[id+1:]...)
	thunder.SetConfig(config)

	err = thunder.SaveConfig()
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	writeJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, fmt.Sprintf(`Source %d deleted!`, id)})
}

func getSourceStats(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, errors.New("ID is not an integer"), "")
		return
	}

	if (id + 1) > len(thunder.GetConfig().Sources) {
		writeJsonError(w, http.StatusBadRequest, errors.New("invalid source"), "")
		return
	}

	source := thunder.GetConfig().Sources[id]
	driver, err := thunder.GetSourceDriver(source.Driver)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err, "")
		return
	}

	driverInstance, err := driver.New(source.Config)
	if err != nil {
		writeJsonError(w, http.StatusUnprocessableEntity, err, "")
		return
	}

	stats, err := driverInstance.Stats()
	if err != nil {
		writeJsonError(w, http.StatusUnprocessableEntity, err, "")
		return
	}
	writeJsonResponse(w, http.StatusOK, stats)
}
