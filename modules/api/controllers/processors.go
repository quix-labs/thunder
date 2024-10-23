package controllers

import (
	"errors"
	"fmt"
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/modules/http_server"
	"github.com/quix-labs/thunder/utils"
	"log"
	"net/http"
	"os"
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
	ID          string              `json:"id"`
	Source      string              `json:"source"` // as source_id
	Targets     []string            `json:"targets"`
	Table       string              `json:"table"`
	PrimaryKeys []string            `json:"primary_keys"`
	Conditions  []thunder.Condition `json:"conditions"`
	Mapping     thunder.JsonMapping `json:"mapping"`
	Index       string              `json:"index"`
	Enabled     bool                `json:"enabled"`
}

func listProcessors(w http.ResponseWriter, r *http.Request) {
	var res []ProcessorApiDetails
	for _, processor := range thunder.Processors.All() {
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
	if err := thunder.Processors.Register("", processor); err != nil {
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
	id := r.PathValue("id")

	var s thunder.JsonProcessor
	err := http_server.DecodeJSONBody(w, r, &s)
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

	if err := thunder.Processors.Update(id, newProcessor); err != nil {
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
	id := r.PathValue("id")
	if err := thunder.Processors.Delete(id); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	if err := thunder.SaveConfig(); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}
	writeJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, fmt.Sprintf(`Processor %d deleted!`, id)})
}

func indexProcessor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	processor, err := thunder.Processors.Get(id)
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
	id := r.PathValue("id")
	processor, err := thunder.Processors.Get(id)
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
	id := r.PathValue("id")
	processor, err := thunder.Processors.Get(id)
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
	id := r.PathValue("id")

	processor, err := thunder.Processors.Get(id)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err, "")
		return
	}

	exporter := r.FormValue("exporter")
	if exporter == "" {
		exporter = "csv"
	}

	filename := r.FormValue("filename")
	if filename == "" {
		filename = fmt.Sprintf("processor-%d.%s", id, exporter)
	}

	exporterInstance, err := thunder.Exporters.Get(exporter)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err, fmt.Sprintf("invalid exporter: %s", exporter))
		return
	}

	// Start indexing
	broadcaster := utils.NewBroadcaster[*thunder.Document, *thunder.Document](func(doc *thunder.Document) *thunder.Document {
		return doc
	})
	broadcaster.Start()
	defer broadcaster.Close()

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "thunder-")
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// Instantiate transformer writer
	if err := exporterInstance.Load(tmpFile); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	// Start receiver
	listenChan, stopListening := broadcaster.NewListenChan()
	transformerErrorChan := make(chan error, 1)
	go func() {
		defer close(transformerErrorChan)
		var position atomic.Uint64
		for doc := range listenChan {
			position.Add(1)
			if position.Load() == 1 {
				if err := exporterInstance.BeforeAll(); err != nil {
					transformerErrorChan <- err
					return
				}
			}

			if err := exporterInstance.WriteDocument(doc, position.Load()); err != nil {
				transformerErrorChan <- err
				return
			}
		}

		if position.Load() >= 1 {
			if err := exporterInstance.AfterAll(); err != nil {
				transformerErrorChan <- err
				return
			}
		}
		transformerErrorChan <- nil
	}()

	sourceErrChan := make(chan error, 1)
	go func() {
		defer close(sourceErrChan)
		sourceErrChan <- processor.Source.Driver.GetDocumentsForProcessor(processor, broadcaster.In(), 0)
	}()

	for {
		select {
		case err := <-sourceErrChan:
			broadcaster.In().Finish()
			stopListening()
			if err != nil {
				writeJsonError(w, http.StatusInternalServerError, err, "")
				return
			}

		case err := <-transformerErrorChan:
			if err != nil {
				writeJsonError(w, http.StatusInternalServerError, err, "")
				broadcaster.In().Finish()
				stopListening()
				return
			}

			// Rewind the temporary file to the beginning
			if _, err := tmpFile.Seek(0, 0); err != nil {
				writeJsonError(w, http.StatusInternalServerError, err, "")
				return
			}

			// Stream temporary file
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
			w.Header().Set("Cache-Control", "private")
			w.Header().Set("Pragma", "private")

			if mimeType := exporterInstance.MimeType(); mimeType != "" {
				w.Header().Set("Content-Type", mimeType)
			}

			http.ServeContent(w, r, filename, time.Now(), tmpFile)
			return
		}
	}
}
