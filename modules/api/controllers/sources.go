package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/modules/http_server"
	"log"
	"net/http"
	"strconv"
)

func SourceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /go-api/sources", ListSources)
	mux.HandleFunc("POST /go-api/sources", CreateSource)
	mux.HandleFunc("DELETE /go-api/sources/{id}", DeleteSource)
	mux.HandleFunc("GET /go-api/sources/{id}/stats", GetSourceStats)
}

func GetSourceStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"ID is not an integer"}`)))
		return
	}

	if (id + 1) > len(thunder.GetConfig().Sources) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"Invalid source"}`)))
		return
	}

	source := thunder.GetConfig().Sources[id]

	driver, err := thunder.GetSourceDriver(source.Driver)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"%s"}`, err)))
		return
	}

	driverInstance, err := driver.New(source.Config)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bytes, _ := json.Marshal(struct {
			Success bool   `json:"success"`
			Error   string `json:"error"`
		}{Success: false, Error: err.Error()})
		w.Write(bytes)
		return
	}

	stats, err := driverInstance.Stats()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		bytes, _ := json.Marshal(struct {
			Success bool   `json:"success"`
			Error   string `json:"error"`
		}{Success: false, Error: err.Error()})
		w.Write(bytes)
		return
	}

	bytes, _ := json.Marshal(stats)
	w.Write(bytes)

}

func DeleteSource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"ID is not an integer"}`)))
		return
	}

	config := thunder.GetConfig()
	config.Sources = append(config.Sources[:id], config.Sources[id+1:]...)
	thunder.SetConfig(config)

	err = thunder.SaveConfig()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"%s"}`, err)))
		return
	}
	w.Write([]byte(fmt.Sprintf(`{"success":true,"message":"Source %d deleted!"}`, id)))

}

func ListSources(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(thunder.GetConfig().Sources)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"%s"}`, err)))
		return
	}

	w.Write(b)
}
func CreateSource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"%s"}`, err)))
		return
	}

	configInstance, err := thunder.ConvertToDynamicConfig(&driver.Config, p.Config)
	if err != nil {
		panic(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	config := thunder.GetConfig()
	config.Sources = append(config.Sources, thunder.Source{Driver: p.Driver, Config: configInstance})
	thunder.SetConfig(config)

	err = thunder.SaveConfig()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"%s"}`, err)))
		return
	}

	w.Write([]byte(`{"success":true,"message":"Source created"}`))
}
