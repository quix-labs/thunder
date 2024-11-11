//go:build !debug

package helpers

import (
	"encoding/json"
	"net/http"
)

func WriteJsonResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		//TODO
	}
}
