//go:build debug

package helpers

import (
	"errors"
	"github.com/quix-labs/thunder"
	"net/http"
	"runtime/debug"
)

type HttpError struct {
	err    error
	msg    string
	status int
}

var nextStatus = http.StatusInternalServerError

func NextCheckStatus(status int) {
	nextStatus = status
}

func CheckErr(err error) {
	defer func() {
		nextStatus = http.StatusInternalServerError
	}()
	if err != nil {
		panic(HttpError{err: err, status: nextStatus})
	}
}

func Must[T any](value T, err error) T {
	defer func() {
		nextStatus = http.StatusInternalServerError
	}()

	msg := ""
	if strValue, ok := any(value).(string); ok {
		msg = strValue
	}

	if err != nil {
		panic(HttpError{err: err, status: nextStatus, msg: msg})
	}
	return value
}

func ErrorMiddleware(next http.Handler, moduleID string) http.HandlerFunc {
	logger := thunder.GetLoggerForModule(moduleID)

	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				statusCode := http.StatusInternalServerError
				msg := ""

				var err error
				switch x := rec.(type) {
				case string:
					err = errors.New(x)
				case HttpError:
					err = x.err
					statusCode = x.status
					msg = x.msg
				case error:
					err = x
				default:
					err = errors.New("unknown panic")
				}

				logger.Error().
					Int("status", statusCode).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("host", r.Host).
					Str("remote_addr", r.RemoteAddr).
					Str("user_agent", r.UserAgent()).
					Stack().Err(err).
					Msg("")

				WriteJsonError(w, statusCode, err, msg)
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
		Stack   string `json:"stack"`
	}{
		Success: false,
		Error:   error.Error(),
		Message: message,
		Stack:   string(debug.Stack()),
	}
	WriteJsonResponse(w, statusCode, payload)
}
