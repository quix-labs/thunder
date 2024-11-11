//go:build debug

package helpers

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

func WriteJsonResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(response); err != nil {
		log.Err(err)
		// TODO
	}
}
