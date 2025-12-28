package handler

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/Vafeda/go_final_project/internal/models"
	"github.com/Vafeda/go_final_project/internal/services"
)

type TaskHandler struct {
	taskService *services.TaskService
}

func NewTaskHandler(taskService *services.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (in *TaskHandler) ReadAll(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	encodeResponse(w, in.taskService.ReadAll(search))
}

func (in *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	if !decodeSuccess(w, r, &task) {
		return
	}

	encodeResponse(w, in.taskService.Create(&task))
}

func (in *TaskHandler) Read(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	// "Если id пустой это жопа"
	fmt.Printf("Тип: %v\n", reflect.TypeOf(in.taskService.Read(id)))
	fmt.Println(in.taskService.Read(id))
	encodeResponse(w, in.taskService.Read(id))
}

func (in *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	if !decodeSuccess(w, r, &task) {
		return
	}

	encodeResponse(w, in.taskService.Update(&task))
}

func (in *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	encodeResponse(w, in.taskService.Delete(id))
}

func (in *TaskHandler) Done(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	encodeResponse(w, in.taskService.Done(id))
}
