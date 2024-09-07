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

func ProcessorRoutes(mux *http.ServeMux) {
	mux.HandleFunc("DELETE /go-api/processors/{id}", deleteProcessor)
	mux.HandleFunc("GET /go-api/processors", listProcessors)
	mux.HandleFunc("POST /go-api/processors", createProcessor)
}

func listProcessors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if len(thunder.GetConfig().Processors) == 0 {
		w.Write([]byte("[]"))
		return
	}

	b, err := json.Marshal(thunder.GetConfig().Processors)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"%s"}`, err)))
		return
	}

	w.Write(b)
}

func createProcessor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var p thunder.Processor

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

	if (p.Source + 1) > len(thunder.GetConfig().Sources) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"Invalid source"}`)))
	}

	config := thunder.GetConfig()
	config.Processors = append(config.Processors, p)
	thunder.SetConfig(config)

	err = thunder.SaveConfig()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"%s"}`, err)))
		return
	}

	w.Write([]byte(`{"success":true,"message":"Processor created"}`))
}

func deleteProcessor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"ID is not an integer"}`)))
	}

	config := thunder.GetConfig()
	config.Processors = append(config.Processors[:id], config.Processors[id+1:]...)
	thunder.SetConfig(config)

	err = thunder.SaveConfig()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"%s"}`, err)))
		return
	}
	w.Write([]byte(fmt.Sprintf(`{"success":true,"message":"Source %d deleted!"}`, id)))
}
