package services

import (
	"fmt"
	"time"

	"github.com/Vafeda/go_final_project/internal/database"
	"github.com/Vafeda/go_final_project/internal/models"
	"github.com/Vafeda/go_final_project/internal/nextdate"
)

type TaskService struct {
	taskDB *database.TaskDB
}

func NewTaskService(taskDB *database.TaskDB) *TaskService {
	return &TaskService{
		taskDB: taskDB,
	}
}

func (in *TaskService) ReadAll(search string) interface{} {
	if search == "" {
		tasks, err := in.taskDB.ReadAll(10)
		if err != nil {
			return models.ErrorResponse{Error: err.Error()}
		}

		return tasks
	}

	date, err := time.Parse("02.01.2006", search)
	if err != nil {
		tasks, err := in.taskDB.ReadByKeyword(search, 10)
		if err != nil {
			return models.ErrorResponse{Error: err.Error()}
		}

		return tasks
	}

	tasks, err := in.taskDB.ReadByDate(date.Format(nextdate.DateFormat), 10)
	if err != nil {
		return models.ErrorResponse{Error: err.Error()}
	}

	return tasks
}

func (in *TaskService) Create(task *models.Task) interface{} {
	if task.Title == "" {
		return models.ErrorResponse{Error: "title empty"}
	}

	if dateIsValid(task.Date) != nil && task.Date != "" {
		return models.ErrorResponse{Error: "fail create task. date invalid"}
	}

	if task.Repeat == "" || task.Date == "" {
		task.Date = time.Now().Format(nextdate.DateFormat)
	} else {
		next, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return models.ErrorResponse{Error: err.Error()}
		}
		task.Date = next
	}

	index, err := in.taskDB.Create(task)
	if err != nil {
		return models.ErrorResponse{Error: err.Error()}
	}

	return models.TaskCreateResponse{ID: index}
}

func (in *TaskService) Read(id string) interface{} {
	task, err := in.taskDB.Read(id)
	fmt.Println(task, err)
	if err != nil {
		return models.ErrorResponse{Error: err.Error()}
	}

	return task
}

func (in *TaskService) Update(task *models.Task) interface{} {
	err := in.taskDB.Update(task)
	if err != nil {
		return models.ErrorResponse{Error: err.Error()}
	}

	return models.EmptyJSONResponse{}
}

func (in *TaskService) Delete(id string) interface{} {
	if id == "" {
		return models.ErrorResponse{Error: "fail delete task. id is empty"}
	}

	err := in.taskDB.Delete(id)
	if err != nil {
		return models.ErrorResponse{Error: err.Error()}
	}

	return models.EmptyJSONResponse{}
}

func (in *TaskService) Done(id string) interface{} {
	response := in.Read(id)
	task, ok := response.(models.Task)
	if !ok {
		return response
	}

	if task.Repeat == "" {
		return in.Delete(id)
	}

	date, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return models.ErrorResponse{Error: err.Error()}
	}

	task.Date = date
	return in.Update(&task)
}

func dateIsValid(date string) error {
	_, err := time.Parse(nextdate.DateFormat, date)
	return err
}
