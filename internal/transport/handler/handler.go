package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Vafeda/go_final_project/internal/nextdate"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
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

	mux.HandleFunc("GET /api/nextdate", nextDayHandler)
	mux.HandleFunc("POST /api/task", taskHandler)
	return mux
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	no := r.FormValue("now")
	now, _ := time.Parse("20060102", no)
	date, _ := nextdate.NextDate(now, r.FormValue("date"), r.FormValue("repeat"))

	fmt.Println(date)
	w.Write([]byte(date))
}

type taskRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	tReq := taskRequest{}
	if err := json.NewDecoder(r.Body).Decode(&tReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//json.NewEncoder(w).Encode(models.OtherError(err))
		return
	}
}
