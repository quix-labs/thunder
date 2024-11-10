package controllers

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/modules/http_server"
	"golang.org/x/sync/errgroup"
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

	mux.HandleFunc("POST /go-api/processors/index", indexAllProcessor)
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
			thunder.GetLoggerForModule("thunder.api").Error().Msg(err.Error())
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
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	// Start source
	var inChan = make(chan *thunder.Document)
	var errChan = make(chan error)

	go func() {
		defer close(inChan)
		errChan <- processor.Source.Driver.GetDocumentsForProcessor(processor, inChan, ctx, 1)
		close(errChan)
	}()

	select {
	case <-ctx.Done():
		writeJsonError(w, http.StatusInternalServerError, ctx.Err(), "")
		return
	case err := <-errChan:
		if err != nil {
			writeJsonError(w, http.StatusInternalServerError, err, "")
		}
	case doc, open := <-inChan:
		if !open {
			// Channel is closed, return 204 No Content
			writeJsonResponse(w, http.StatusNoContent, nil)
			return
		}

		// Successfully received a document, write it to the response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(doc.Json)
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
			thunder.GetLoggerForModule("thunder.api").Error().Msg(err.Error())
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
			thunder.GetLoggerForModule("thunder.api").Error().Msg(err.Error())
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
	}{true, fmt.Sprintf("Processor %s updated", id)})
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
	}{true, fmt.Sprintf(`Processor %s deleted!`, id)})
}

func indexAllProcessor(w http.ResponseWriter, r *http.Request) {

	sync := r.URL.Query().Has("sync")

	if sync {
		var eg, egCtx = errgroup.WithContext(r.Context())

		processors := thunder.Processors.All()

		groupedErr := make(map[string]string, len(processors))
		for _, processor := range processors {
			//id := id
			eg.Go(func() error {
				start := time.Now()
				err := processor.FullIndex(egCtx)
				if err != nil {
					groupedErr[processor.Index] = err.Error()
				} else {
					groupedErr[processor.Index] = fmt.Sprintf("indexed, took %fs!", time.Since(start).Seconds())
				}
				return nil // Do not return error but append to error array
			})
		}

		if err := eg.Wait(); err != nil {
			writeJsonError(w, http.StatusInternalServerError, err, "")
			return
		}

		writeJsonResponse(w, http.StatusOK, struct {
			Success bool              `json:"success"`
			Results map[string]string `json:"results"`
		}{true, groupedErr})
		return
	}

	// BACKGROUND RUN
	for id, processor := range thunder.Processors.All() {
		go func() {
			thunder.GetLoggerForModule("thunder.api").Info().Msgf("Processor %s has started in the background", id)
			err := processor.FullIndex(context.Background())
			if err != nil {
				thunder.GetLoggerForModule("thunder.api").Error().Msgf("an error occurred during indexing of processor %s: %s", id, err.Error())
				return
			}
			thunder.GetLoggerForModule("thunder.api").Info().Msgf("Processor %s has finished in the background", id)
		}()
	}

	writeJsonResponse(w, http.StatusAccepted, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, "Indexing claimed"})
}

func indexProcessor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	processor, err := thunder.Processors.Get(id)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err, "")
		return
	}

	sync := r.URL.Query().Has("sync")

	if sync {
		if err := processor.FullIndex(r.Context()); err != nil {
			writeJsonError(w, http.StatusInternalServerError, err, "error during indexing")
			return
		}
		writeJsonResponse(w, http.StatusOK, struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{true, "indexing succeed"})
		return
	}

	go func() {
		thunder.GetLoggerForModule("thunder.api").Info().Msgf("Processor %s has started in the background", id)
		err = processor.FullIndex(context.Background())
		if err != nil {
			thunder.GetLoggerForModule("thunder.api").Error().Msgf("an error occurred during indexing of processor %s: %s", id, err.Error())
			return
		}
		thunder.GetLoggerForModule("thunder.api").Info().Msgf("Processor %s has finished in the background", id)
	}()

	writeJsonResponse(w, http.StatusAccepted, struct {
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
	}{true, fmt.Sprintf("Processor %s started", id)})
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
	}{true, fmt.Sprintf("Processor %s stopped", id)})
}

func downloadProcessor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	processor, err := thunder.Processors.Get(id)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err, "")
		return
	}

	exporter := cmp.Or(r.FormValue("exporter"), "thunder.csv")
	filename := cmp.Or(r.FormValue("filename"), fmt.Sprintf("processor-%s.%s", id, exporter))

	exporterInstance, err := thunder.Exporters.Get(exporter)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err, fmt.Sprintf("invalid exporter: %s", exporter))
		return
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "thunder-")
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	defer func(name string) {
		if err := os.Remove(name); err != nil {
			thunder.GetLoggerForModule("thunder.api").Error().Msg(err.Error())
			return
		}
		thunder.GetLoggerForModule("thunder.api").Debug().Msgf("temporary resource %s removed", tmpFile.Name())
	}(tmpFile.Name())

	defer func(tmpFile *os.File) {
		if err := tmpFile.Close(); err != nil {
			thunder.GetLoggerForModule("thunder.api").Error().Msg(err.Error())
			return
		}
		thunder.GetLoggerForModule("thunder.api").Debug().Msgf("temporary resource %s closed", tmpFile.Name())
	}(tmpFile)

	// Instantiate transformer writer
	if err := exporterInstance.Load(tmpFile); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}

	// Start indexing
	eg, egCtx := errgroup.WithContext(r.Context())

	// Start source
	var inChan = make(chan *thunder.Document)
	eg.Go(func() error {
		defer close(inChan)
		return processor.Source.Driver.GetDocumentsForProcessor(processor, inChan, egCtx, 0)
	})

	// Start exporter broadcasting
	eg.Go(func() error {
		var position atomic.Uint64
		for {
			select {
			case <-egCtx.Done():
				return egCtx.Err()
			case <-time.After(time.Second * 10):
				return context.DeadlineExceeded
			case doc, open := <-inChan:
				if !open {
					if position.Load() >= 1 {
						return exporterInstance.AfterAll()
					}
					return nil
				}
				position.Add(1)
				if position.Load() == 1 {
					if err := exporterInstance.BeforeAll(); err != nil {
						return err
					}
				}

				if err := exporterInstance.WriteDocument(doc, position.Load()); err != nil {
					return err
				}
			}
		}
	})

	err = eg.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			err = context.Cause(egCtx)
		}
		writeJsonError(w, http.StatusInternalServerError, err, "")
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
