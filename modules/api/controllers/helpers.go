package controllers

import (
	"encoding/json"
	"github.com/quix-labs/thunder"
	"log"
	"net/http"
)

func writeJsonResponse(w http.ResponseWriter, statusCode int, response interface{}) {

	marshal, err := json.Marshal(response)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(marshal)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err, "")
	}
}

func writeJsonError(w http.ResponseWriter, statusCode int, error error, message string) {
	var payload = struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
		Message string `json:"message"`
	}{
		Success: false,
		Error:   error.Error(),
		Message: message,
	}
	marshal, err := json.Marshal(&payload)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(marshal)
	if err != nil {
		log.Println(err)
	}
}

func GetProcessorStatus() map[int]thunder.ProcessorStatus {
	processorsStates := thunder.GetProcessors()
	var response = make(map[int]thunder.ProcessorStatus, len(processorsStates))
	for idx, processor := range processorsStates {
		response[idx] = processor.Status
	}
	return response
}
