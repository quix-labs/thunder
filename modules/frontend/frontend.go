package frontend

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

func Start() error {
	fmt.Println("Server listening on port :3000")
	fmt.Println("Access http://localhost:3000")

	srv := &http.Server{
		Addr:    ":3000",
		Handler: router(),
	}

	return srv.ListenAndServe()
}

//go:embed all:build
var StaticFiles embed.FS

func router() http.Handler {
	mux := http.NewServeMux()

	// static files
	staticFS, _ := fs.Sub(StaticFiles, "build")
	mux.Handle("/", FileServerWith404(http.FS(staticFS), func(w http.ResponseWriter, r *http.Request) (doDefaultFileServe bool) {
		if strings.HasPrefix(r.URL.Path, "/go-api/") {
			return true
		}

		errorFile, _ := staticFS.Open("404.html")
		defer errorFile.Close()

		fileInfo, err := errorFile.Stat()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Read the file content
		content, err := io.ReadAll(errorFile)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(404)

		http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), bytes.NewReader(content))
		return
	}))

	SourceRoutes(mux)
	SourceDriverRoutes(mux)
	ProcessorRoutes(mux)

	return mux
}

type FSHandler404 = func(w http.ResponseWriter, r *http.Request) (doDefaultFileServe bool)

func FileServerWith404(root http.FileSystem, handler404 FSHandler404) http.Handler {
	fs := http.FileServer(root)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upath := r.URL.Path
		if !strings.HasPrefix(upath, "/") {
			upath = "/" + upath
			r.URL.Path = upath
		}
		upath = path.Clean(upath)

		f, err := root.Open(upath)
		if err != nil {
			if os.IsNotExist(err) {
				if handler404 != nil {
					doDefault := handler404(w, r)
					if !doDefault {
						return
					}
				}
			}
		}

		if err == nil {
			f.Close()
		}

		fs.ServeHTTP(w, r)
	})
}
