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
	"time"
)

func ProcessorRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /go-api/processors", listProcessors)
	mux.HandleFunc("POST /go-api/processors/test", testProcessor)

	mux.HandleFunc("POST /go-api/processors", createProcessor)
	mux.HandleFunc("PUT /go-api/processors/{id}", updateProcessor)
	mux.HandleFunc("DELETE /go-api/processors/{id}", deleteProcessor)

}

func listProcessors(w http.ResponseWriter, r *http.Request) {
	writeJsonResponse(w, http.StatusOK, thunder.GetConfig().Processors)
}

func testProcessor(w http.ResponseWriter, r *http.Request) {
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
		writeJsonError(w, http.StatusBadRequest, errors.New("invalid source"), "")
		return
	}
	for _, targetId := range p.Targets {
		if (targetId + 1) > len(thunder.GetConfig().Targets) {
			writeJsonError(w, http.StatusBadRequest, errors.New(fmt.Sprintf("invalid target %d", targetId)), "")
			return
		}
	}
	// Get related
	source := thunder.GetConfig().Sources[p.Source]
	driverInfo, err := thunder.GetSourceDriver(source.Driver)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err, "")
		return
	}

	driver, err := driverInfo.New(source.Config)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}
	defer driver.Shutdown()

	// TRY TO FETCH 1 DOCUMENT
	docChan, errChan := driver.GetDocumentsForProcessor(&p, 1)
	select {
	case doc, open := <-docChan:
		if !open {
			writeJsonResponse(w, http.StatusOK, nil)
			return
		}

		keyBytes, err := json.Marshal(doc.PrimaryKeys)
		if err != nil {
			writeJsonError(w, http.StatusInternalServerError, err, "")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"primary_keys":%s,"document":%s}`, keyBytes, doc.Json)))
		return
	case err := <-errChan:
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	case <-time.After(time.Second * 5):
		writeJsonError(w, http.StatusInternalServerError, err, "timeout reached")
		return
	}
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
		return
	}
	for _, targetId := range p.Targets {
		if (targetId + 1) > len(thunder.GetConfig().Targets) {
			writeJsonError(w, http.StatusBadRequest, errors.New(fmt.Sprintf("invalid target %d", targetId)), "")
			return
		}
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

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"success":true,"message":"Processor created"}`))
}

func updateProcessor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"ID is not an integer"}`)))
		return
	}

	var p thunder.Processor

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

	if (p.Source + 1) > len(thunder.GetConfig().Sources) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"Invalid source"}`)))
		return
	}

	for _, targetId := range p.Targets {
		if (targetId + 1) > len(thunder.GetConfig().Targets) {
			writeJsonError(w, http.StatusBadRequest, errors.New(fmt.Sprintf("invalid target %d", targetId)), "")
			return
		}
	}

	config := thunder.GetConfig()
	config.Processors[id] = p
	thunder.SetConfig(config)

	err = thunder.SaveConfig()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"%s"}`, err)))
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"success":true,"message":"Processor %d updated!"}`, id)))
}

func deleteProcessor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"success":false,"error":"ID is not an integer"}`)))
		return
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
	w.Write([]byte(fmt.Sprintf(`{"success":true,"message":"Processor %d deleted!"}`, id)))
}
