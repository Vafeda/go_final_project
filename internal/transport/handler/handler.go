package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Vafeda/go_final_project/internal/database"
	"github.com/Vafeda/go_final_project/internal/services"
)

func SetupRoutes(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	web, err := NewWeb()
	if err != nil {
		return nil
	}
	mux.HandleFunc("/", web.Get)

	mux.HandleFunc("GET /api/nextdate", nextDayHandler)

	mux.HandleFunc("POST /api/signin", signIn)

	taskHandler := NewTaskHandler(services.NewTaskService(database.NewTaskDB(db)))
	mux.HandleFunc("POST /api/task", auth(taskHandler.Create))
	mux.HandleFunc("GET /api/task", auth(taskHandler.Read))
	mux.HandleFunc("PUT /api/task", auth(taskHandler.Update))
	mux.HandleFunc("DELETE /api/task", auth(taskHandler.Delete))
	mux.HandleFunc("GET /api/tasks", auth(taskHandler.ReadAll))
	mux.HandleFunc("POST /api/task/done", auth(taskHandler.Done))

	return mux
}

func encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}
