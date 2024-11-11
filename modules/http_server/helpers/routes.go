package helpers

import "net/http"

func PrefixRoutes(prefix string, mux *http.ServeMux) *http.ServeMux {
	prefixedMux := http.NewServeMux()
	prefixedMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = prefix + r.URL.Path
		mux.ServeHTTP(w, r)
	})
	return prefixedMux
}
