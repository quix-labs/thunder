package frontend

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/quix-labs/thunder"
	"log"
	"net/http"
	"strconv"
)

func SourceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /go-api/sources", listSources)
	mux.HandleFunc("POST /go-api/sources", createSource)
	mux.HandleFunc("DELETE /go-api/sources/{id}", deleteSource)

	mux.HandleFunc("GET /go-api/sources/{id}/stats", getSourceStats)
}

func getSourceStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"ID is not an integer"}`)))
	}

	if (id + 1) > len(thunder.GetConfig().Sources) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"Invalid source"}`)))
	}

	source := thunder.GetConfig().Sources[id]

	driver, err := thunder.GetSourceDriver(source.Driver)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"%s"}`, err)))
	}

	stats, err := driver.New().Stats(source.Config)
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

func deleteSource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"ID is not an integer"}`)))
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

func listSources(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(thunder.GetConfig().Sources)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"%s"}`, err)))
		return
	}

	w.Write(b)
}
func createSource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var p struct {
		Driver string `json:"driver"`
		Config any    `json:"config"`
	}
	err := decodeJSONBody(w, r, &p)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
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

	configInstance, err := thunder.ConvertSourceDriverConfig(&driver, p.Config)
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
