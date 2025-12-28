package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

const webDir = "web"

type WebHandler struct {
	fs http.Handler
}

func NewWeb() (*WebHandler, error) {
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("folder at path \"%s\" not found", webDir)
	}
	return &WebHandler{fs: http.FileServer(http.Dir(webDir))}, nil
}

func (wH *WebHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
		return
	}

	wH.fs.ServeHTTP(w, r)
}
