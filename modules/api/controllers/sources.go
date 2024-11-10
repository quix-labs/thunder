package controllers

import (
	"errors"
	"fmt"
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/modules/http_server"
	"net/http"
)

func SourceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /go-api/sources", listSources)
	mux.HandleFunc("POST /go-api/sources", createSource)
	mux.HandleFunc("PUT /go-api/sources/{id}", updateSource)
	mux.HandleFunc("DELETE /go-api/sources/{id}", deleteSource)
	mux.HandleFunc("GET /go-api/sources/{id}/stats", getSourceStats)
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
		serializeSource, err := thunder.SerializeSource(&source)
		if err != nil {
			writeJsonError(w, http.StatusInternalServerError, err, "")
			return
		}
		res = append(res, SourceApiDetails{
			ID:      serializeSource.ID,
			Driver:  serializeSource.Driver,
			Config:  serializeSource.Config,
			Excerpt: source.Config.Excerpt(),
		})
	}
	writeJsonResponse(w, http.StatusOK, res)
}

func createSource(w http.ResponseWriter, r *http.Request) {
	var p thunder.JsonSource

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

	source, err := thunder.UnserializeSource(&p)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err, "")
		return
	}

	err = thunder.Sources.Register("", *source)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

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
	id := r.PathValue("id")

	var s thunder.JsonSource
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

	newSource, err := thunder.UnserializeSource(&s)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	if err := thunder.Sources.Update(id, *newSource); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	if err = thunder.SaveConfig(); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	writeJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, fmt.Sprintf("Source %d updated", id)})
}

func deleteSource(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := thunder.Sources.Delete(id)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

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
	id := r.PathValue("id")

	source, err := thunder.Sources.Get(id)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err, "")
		return
	}

	stats, err := source.Driver.Stats()
	if err != nil {
		writeJsonError(w, http.StatusUnprocessableEntity, err, "")
		return
	}

	writeJsonResponse(w, http.StatusOK, stats)
}
