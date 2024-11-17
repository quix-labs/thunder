package frontend

//go:generate sh -c "cd src && yarn"
//go:generate sh -c "cd src && NUXT_OUTPUT_DIR=../static yarn generate"

import (
	"bytes"
	"embed"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

//go:embed all:static
var StaticFiles embed.FS

func HandleFrontend(mux *http.ServeMux) {
	// static files
	staticFS, _ := fs.Sub(StaticFiles, "static")
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

		// CANNOT RETURN 404 due to SPA
		//w.WriteHeader(404)

		http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), bytes.NewReader(content))
		return
	}))
}

type FSHandler404 = func(w http.ResponseWriter, r *http.Request) (doDefaultFileServe bool)

func FileServerWith404(root http.FileSystem, handler404 FSHandler404) http.Handler {
	fileServer := http.FileServer(root)

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

		fileServer.ServeHTTP(w, r)
	})
}
