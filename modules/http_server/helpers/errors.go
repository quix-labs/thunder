//go:build !debug

package helpers

import (
	"errors"
	"github.com/quix-labs/thunder"
	"net/http"
)

func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ErrorMiddleware(next http.Handler, moduleID string) http.HandlerFunc {
	logger := thunder.GetLoggerForModule(moduleID)

	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				var err error
				switch x := rec.(type) {
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					err = errors.New("unknown panic")
				}

				logger.Error().
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("host", r.Host).
					Str("remote_addr", r.RemoteAddr).
					Str("user_agent", r.UserAgent()).
					Err(err).
					Msg("")

				WriteJsonError(w, http.StatusInternalServerError, err, "")
			}
		}()
		next.ServeHTTP(w, r)
	}

}

func WriteJsonError(w http.ResponseWriter, statusCode int, error error, message string) {
	var payload = struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
		Message string `json:"message"`
	}{
		Success: false,
		Error:   error.Error(),
		Message: message,
	}

	WriteJsonResponse(w, statusCode, payload)
}
