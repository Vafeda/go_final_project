package handler

import (
	"database/sql"
	"encoding/json"
	"github.com/Vafeda/go_final_project/internal/database"
	"github.com/Vafeda/go_final_project/internal/services"
	"net/http"
	"time"

	"github.com/Vafeda/go_final_project/internal/models"
	"github.com/Vafeda/go_final_project/internal/nextdate"
)

func SetupRoutes(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	// Проверяем существование папки
	web, err := NewWeb()
	if err != nil {
		return nil
	}
	mux.HandleFunc("/", web.Get)

	mux.HandleFunc("GET /api/nextdate", nextDayHandler)
	//
	taskHandler := NewTaskHandler(services.NewTaskService(database.NewTaskDB(db)))
	mux.HandleFunc("GET /api/tasks", auth(taskHandler.ReadAll))
	//
	mux.HandleFunc("POST /api/task", auth(taskHandler.Create))
	mux.HandleFunc("GET /api/task", auth(taskHandler.Read))
	mux.HandleFunc("PUT /api/task", auth(taskHandler.Update))
	mux.HandleFunc("DELETE /api/task", auth(taskHandler.Delete))
	//
	mux.HandleFunc("POST /api/task/done", auth(taskHandler.Done))

	//
	mux.HandleFunc("POST /api/signin", signIn)
	return mux
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse("20060102", r.FormValue("now"))
	if err != nil {
		encodeResponse(w, models.ErrorResponse{})
	}

	date, _ := nextdate.NextDate(now, r.FormValue("date"), r.FormValue("repeat"))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(date))
}

func decodeSuccess(w http.ResponseWriter, r *http.Request, request interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		encodeResponse(w, &models.ErrorResponse{Error: err.Error()})
		return false
	}
	return true
}

func encodeResponse(w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json")
	switch response.(type) {
	case models.TaskCreateResponse:

		w.WriteHeader(http.StatusCreated)
	case models.Task:
		w.WriteHeader(http.StatusOK)
	case models.Tasks:
		w.WriteHeader(http.StatusOK)
	case models.ErrorResponse:
		w.WriteHeader(http.StatusBadRequest)
	case models.EmptyJSONResponse:
		w.WriteHeader(http.StatusOK)
	case models.JSONToken:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	// Возможно стоит просто залогировать ошибку
	json.NewEncoder(w).Encode(response)
}
