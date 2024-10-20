package controllers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/modules/http_server"
	"github.com/quix-labs/thunder/utils"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

func ProcessorRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /go-api/processors", listProcessors)
	mux.HandleFunc("POST /go-api/processors/test", testProcessor)

	mux.HandleFunc("POST /go-api/processors", createProcessor)
	mux.HandleFunc("PUT /go-api/processors/{id}", updateProcessor)
	mux.HandleFunc("DELETE /go-api/processors/{id}", deleteProcessor)

	mux.HandleFunc("POST /go-api/processors/{id}/index", indexProcessor)
	mux.HandleFunc("POST /go-api/processors/{id}/start", startProcessor)
	mux.HandleFunc("POST /go-api/processors/{id}/stop", stopProcessor)

	mux.HandleFunc("GET /go-api/processors/{id}/download", downloadProcessor)

}

type ProcessorApiDetails struct {
	Indexing    bool                `json:"indexing"`
	Listening   bool                `json:"listening"`
	ID          int                 `json:"id"`
	Source      int                 `json:"source"`  // as source_id
	Targets     []int               `json:"targets"` // as targets_id
	Table       string              `json:"table"`
	PrimaryKeys []string            `json:"primary_keys"`
	Conditions  []thunder.Condition `json:"conditions"`
	Mapping     thunder.JsonMapping `json:"mapping"`
	Index       string              `json:"index"`
	Enabled     bool                `json:"enabled"`
}

func listProcessors(w http.ResponseWriter, r *http.Request) {
	processors := thunder.GetProcessors()
	var res []ProcessorApiDetails
	for _, processor := range processors {
		serializedProcessor, err := thunder.SerializeProcessor(processor)
		if err != nil {
			writeJsonError(w, http.StatusInternalServerError, err, "")
			return
		}
		res = append(res, ProcessorApiDetails{
			Indexing:  processor.Indexing.Load(),
			Listening: processor.Listening.Load(),

			ID:          serializedProcessor.ID,
			Source:      serializedProcessor.Source,
			Targets:     serializedProcessor.Targets,
			Table:       serializedProcessor.Table,
			PrimaryKeys: serializedProcessor.PrimaryKeys,
			Conditions:  serializedProcessor.Conditions,
			Mapping:     serializedProcessor.Mapping,
			Index:       serializedProcessor.Index,
			Enabled:     serializedProcessor.Enabled,
		})
	}
	writeJsonResponse(w, http.StatusOK, res)
}

func testProcessor(w http.ResponseWriter, r *http.Request) {
	var p thunder.JsonProcessor
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

	processor, err := thunder.UnserializeProcessor(&p)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	// Start indexing
	broadcaster := utils.NewBroadcaster[*thunder.Document, *thunder.Document](func(doc *thunder.Document) *thunder.Document {
		return doc
	})
	broadcaster.Start()
	defer broadcaster.Close()

	var errChan = make(chan error)

	// Start source
	go func() {
		defer broadcaster.In().Finish()
		if err := processor.Source.Driver.GetDocumentsForProcessor(processor, broadcaster.In(), 1); err != nil {
			errChan <- err
		}
	}()

	// Start receiver
	listenChan, _ := broadcaster.NewListenChan()

	select {
	case doc, open := <-listenChan:
		if !open {
			writeJsonResponse(w, http.StatusOK, nil)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(doc.Json)
		return
	case err := <-errChan:
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	case <-time.After(time.Second * 5):
		writeJsonError(w, http.StatusInternalServerError, errors.New("timeout reached"), "database took more than 5s to generate document")
		return
	}
}

func createProcessor(w http.ResponseWriter, r *http.Request) {
	var p thunder.JsonProcessor

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

	processor, err := thunder.UnserializeProcessor(&p)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err, "")
		return
	}
	processor.ID = 0
	if err := thunder.AddProcessor(processor); err != nil {
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
	}{true, "Processor created"})
}

func updateProcessor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, errors.New("ID is not an integer"), "")
		return
	}

	var s thunder.JsonProcessor
	err = http_server.DecodeJSONBody(w, r, &s)
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

	newProcessor, err := thunder.UnserializeProcessor(&s)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	if err := thunder.UpdateProcessor(id, newProcessor); err != nil {
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
	}{true, fmt.Sprintf("Processor %d updated", id)})
}

func deleteProcessor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, errors.New("ID is not an integer"), "")
		return
	}

	err = thunder.DeleteProcessor(id)
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
	}{true, fmt.Sprintf(`Processor %d deleted!`, id)})
}

func indexProcessor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, errors.New("ID is not an integer"), "")
		return
	}

	processor, err := thunder.GetProcessor(id)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err, "")
		return
	}

	go func() {
		if err := processor.FullIndex(); err != nil {
			return // TODO
		}
	}()

	writeJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, "Indexing claimed"})
}

func startProcessor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, errors.New("ID is not an integer"), "")
		return
	}

	processor, err := thunder.GetProcessor(id)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err, "")
		return
	}

	go func() {
		if err := processor.Start(); err != nil {
			thunder.GetLoggerForProcessor(processor).Error().Msg(err.Error())
			return
		}
	}()

	writeJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, fmt.Sprintf("Processor %d started", id)})
}

func stopProcessor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, errors.New("ID is not an integer"), "")
		return
	}

	processor, err := thunder.GetProcessor(id)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err, "")
		return
	}

	if err := processor.Stop(); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
	}

	writeJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, fmt.Sprintf("Processor %d stopped", id)})
}

func downloadProcessor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, errors.New("ID is not an integer"), "")
		return
	}

	processor, err := thunder.GetProcessor(id)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err, "")
		return
	}

	format := r.FormValue("format")
	if format == "" {
		format = "csv"
	}

	filename := r.FormValue("filename")
	if filename == "" {
		filename = fmt.Sprintf("processor-%d.%s", id, format)
	}

	// Send download file header
	w.Header().Set("content-disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("content-transfer-encoding", "binary")
	w.Header().Set("transfer-encoding", "chunked")
	w.Header().Set("accept-ranges", "bytes")
	w.Header().Set("cache-control", "private")
	w.Header().Set("pragma", "private")

	// Start indexing
	broadcaster := utils.NewBroadcaster[*thunder.Document, *thunder.Document](func(doc *thunder.Document) *thunder.Document {
		return doc
	})
	broadcaster.Start()
	defer broadcaster.Close()

	// Start receiver
	listenChan, stopListening := broadcaster.NewListenChan()
	writed := make(chan struct{})
	switch format {
	case "json":
		go func() {
			w.Header().Set("content-type", "application/json")
			w.Write([]byte("["))
			var i atomic.Uint64
			for doc := range listenChan {
				if i.Load() > 0 {
					w.Write([]byte(","))
				}

				w.Write(doc.Json)
			}
			w.Write([]byte("]"))
			writed <- struct{}{}
		}()
		break
	case "csv":
		go func() {
			w.Header().Set("content-type", "text/csv")

			csvWriter := csv.NewWriter(w)
			csvWriter.Write([]string{"pkey", "json"})

			for doc := range listenChan {
				csvWriter.Write([]string{doc.Pkey, string(doc.Json)})
			}

			writed <- struct{}{}
		}()
		break
	default:
		writeJsonError(w, http.StatusBadRequest, errors.New(fmt.Sprintf("unsupported format: %s", format)), "")
		return
	}

	// Start source
	if err := processor.Source.Driver.GetDocumentsForProcessor(processor, broadcaster.In(), 0); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
	}

	broadcaster.In().Finish()
	stopListening()

	<-writed
}
