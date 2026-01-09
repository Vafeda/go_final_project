package handler

import (
	"fmt"
	"net/http"

	"github.com/Vafeda/TODO-List/internal/models"
	"github.com/Vafeda/TODO-List/internal/services"
)

const (
	keySearch = "search"
	keyId     = "id"
)

type TaskHandler struct {
	taskService *services.TaskService
}

func NewTaskHandler(taskService *services.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (in *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	task, err := decode[models.Task](r)
	if err != nil {
		err = encode(w, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	taskResponse, err := in.taskService.Create(&task)
	fmt.Println(taskResponse, err)
	if err != nil {
		err = encode(w, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	err = encode(w, http.StatusCreated, taskResponse)
}

func (in *TaskHandler) Read(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue(keyId)

	task, err := in.taskService.Read(id)
	if err != nil {
		err = encode(w, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	err = encode(w, http.StatusOK, task)
}

func (in *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	task, err := decode[models.Task](r)
	if err != nil {
		err = encode(w, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	err = in.taskService.Update(&task)
	if err != nil {
		err = encode(w, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	err = encode(w, http.StatusOK, models.EmptyJSONResponse{})
}

func (in *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue(keyId)

	err := in.taskService.Delete(id)
	if err != nil {
		err = encode(w, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	err = encode(w, http.StatusOK, models.EmptyJSONResponse{})
}

func (in *TaskHandler) ReadAll(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue(keySearch)

	tasks, err := in.taskService.ReadAll(search)
	if err != nil {
		err = encode(w, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	err = encode(w, http.StatusOK, tasks)
}

func (in *TaskHandler) Done(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue(keyId)

	err := in.taskService.Done(id)
	if err != nil {
		err = encode(w, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	err = encode(w, http.StatusOK, models.EmptyJSONResponse{})
}
