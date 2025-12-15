package handler

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	webPath := "D:\\Golang\\go_final_project\\web"

	// Проверяем существование папки
	if _, err := os.Stat(webPath); os.IsNotExist(err) {
		log.Fatalf("Папка не найдена: %s", webPath)
	}

	fs := http.FileServer(http.Dir(webPath))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			http.ServeFile(w, r, filepath.Join(webPath, "index.html"))
			return
		}
		
		fs.ServeHTTP(w, r)
	})

	return mux
}
