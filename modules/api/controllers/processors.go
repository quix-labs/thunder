package controllers

import (
	"cmp"
	"context"
	"fmt"
	"github.com/quix-labs/thunder"
	"github.com/quix-labs/thunder/modules/http_server"
	"github.com/quix-labs/thunder/modules/http_server/helpers"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"strconv"
	"time"
)

func ProcessorRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /processors", listProcessors)
	mux.HandleFunc("POST /processors/test", testProcessor)

	mux.HandleFunc("POST /processors", createProcessor)
	mux.HandleFunc("PUT /processors/{id}", updateProcessor)
	mux.HandleFunc("DELETE /processors/{id}", deleteProcessor)

	mux.HandleFunc("POST /processors/index", indexAllProcessor)
	mux.HandleFunc("POST /processors/{id}/index", indexProcessor)
	mux.HandleFunc("POST /processors/{id}/start", startProcessor)
	mux.HandleFunc("POST /processors/{id}/stop", stopProcessor)

	mux.HandleFunc("GET /processors/{id}/download", downloadProcessor)

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
		serializedProcessor := helpers.Must(thunder.SerializeProcessor(processor))
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
	helpers.WriteJsonResponse(w, http.StatusOK, res)
}

func testProcessor(w http.ResponseWriter, r *http.Request) {
	var p thunder.JsonProcessor
	helpers.CheckErr(http_server.DecodeJSONBody(w, r, &p))

	processor := helpers.Must(thunder.UnserializeProcessor(&p))

	// Start indexing
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	// Start source
	var inChan = make(chan *thunder.Document)
	var errChan = make(chan error)

	go func() {
		defer close(inChan)
		errChan <- processor.Source.Driver.GetDocumentsForProcessor(processor, inChan, ctx, 1)
	}()

	select {
	case <-ctx.Done():
		helpers.CheckErr(ctx.Err())
		return
	case err := <-errChan:
		close(errChan)
		helpers.CheckErr(err)
	case doc, open := <-inChan:
		if !open {
			// Channel is closed, return 204 No Content
			helpers.WriteJsonResponse(w, http.StatusNoContent, nil)
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
	helpers.CheckErr(http_server.DecodeJSONBody(w, r, &p))

	helpers.NextCheckStatus(http.StatusBadRequest)
	processor := helpers.Must(thunder.UnserializeProcessor(&p))

	helpers.CheckErr(thunder.Processors.Register("", processor))
	helpers.CheckErr(thunder.SaveConfig())

	helpers.WriteJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, "Processor created"})
}

func updateProcessor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var s thunder.JsonProcessor
	helpers.CheckErr(http_server.DecodeJSONBody(w, r, &s))

	newProcessor := helpers.Must(thunder.UnserializeProcessor(&s))

	helpers.CheckErr(thunder.Processors.Update(id, newProcessor))

	helpers.CheckErr(thunder.SaveConfig())

	helpers.WriteJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, fmt.Sprintf("Processor %s updated", id)})
}

func deleteProcessor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	helpers.CheckErr(thunder.Processors.Delete(id))
	helpers.CheckErr(thunder.SaveConfig())

	helpers.WriteJsonResponse(w, http.StatusOK, struct {
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

		helpers.CheckErr(eg.Wait())

		helpers.WriteJsonResponse(w, http.StatusOK, struct {
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

	helpers.WriteJsonResponse(w, http.StatusAccepted, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, "Indexing claimed"})
}

func indexProcessor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	helpers.NextCheckStatus(http.StatusBadRequest)
	processor := helpers.Must(thunder.Processors.Get(id))

	sync := r.URL.Query().Has("sync")

	if sync {
		helpers.CheckErr(processor.FullIndex(r.Context()))

		helpers.WriteJsonResponse(w, http.StatusOK, struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{true, "indexing succeed"})
		return
	}

	go func() {
		thunder.GetLoggerForModule("thunder.api").Info().Msgf("Processor %s has started in the background", id)
		if err := processor.FullIndex(context.Background()); err != nil {
			thunder.GetLoggerForModule("thunder.api").Error().Msgf("an error occurred during indexing of processor %s: %s", id, err.Error())
			return
		}
		thunder.GetLoggerForModule("thunder.api").Info().Msgf("Processor %s has finished in the background", id)
	}()

	helpers.WriteJsonResponse(w, http.StatusAccepted, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, "Indexing claimed"})
}

func startProcessor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	helpers.NextCheckStatus(http.StatusBadRequest)
	processor := helpers.Must(thunder.Processors.Get(id))

	go func() {
		if err := processor.Start(); err != nil {
			thunder.GetLoggerForProcessor(processor).Error().Msg(err.Error())
		}
	}()

	helpers.WriteJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, fmt.Sprintf("Processor %s started", id)})
}

func stopProcessor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	helpers.NextCheckStatus(http.StatusBadRequest)
	processor := helpers.Must(thunder.Processors.Get(id))

	helpers.CheckErr(processor.Stop())

	helpers.WriteJsonResponse(w, http.StatusOK, struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{true, fmt.Sprintf("Processor %s stopped", id)})
}

func downloadProcessor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	helpers.NextCheckStatus(http.StatusBadRequest)
	processor := helpers.Must(thunder.Processors.Get(id))

	exporter := cmp.Or(r.FormValue("exporter"), "thunder.csv")
	exporterInstance := helpers.Must(thunder.Exporters.Get(exporter))

	// Create temporary file
	tmpFile := helpers.Must(os.CreateTemp("", "thunder-"))

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

	// Start the streaming process
	var limit uint64 = 0
	if r.FormValue("limit") != "" {
		limit = helpers.Must(strconv.ParseUint(r.FormValue("limit"), 10, 64))
	}

	helpers.CheckErr(processor.StreamDocuments(r.Context(), exporterInstance, tmpFile, limit))

	// Rewind the temporary file to the beginning
	_ = helpers.Must(tmpFile.Seek(0, 0))

	// Stream temporary file
	filename := cmp.Or(r.FormValue("filename"), fmt.Sprintf("processor-%s.%s", id, exporter))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Cache-Control", "private")
	w.Header().Set("Pragma", "private")

	if mimeType := exporterInstance.MimeType(); mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
	}

	http.ServeContent(w, r, filename, time.Now(), tmpFile)
	return
}
