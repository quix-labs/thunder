package controllers

import (
	"encoding/json"
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
	_, err = w.Write(append(marshal, '\n'))
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

	_, err = w.Write(append(marshal, '\n'))
	if err != nil {
		log.Println(err)
	}
}
