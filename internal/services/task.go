package services

import (
	"fmt"
	"time"

	"github.com/Vafeda/TODO-List/internal/database"
	"github.com/Vafeda/TODO-List/internal/models"
	"github.com/Vafeda/TODO-List/internal/nextdate"
)

type TaskService struct {
	taskDB *database.TaskDB
}

func NewTaskService(taskDB *database.TaskDB) *TaskService {
	return &TaskService{
		taskDB: taskDB,
	}
}

func (in *TaskService) ReadAll(search string) (models.Tasks, error) {
	if search == "" {
		tasks, err := in.taskDB.ReadAll(10)
		if err != nil {
			return models.Tasks{}, err
		}

		return tasks, nil
	}

	date, err := time.Parse("02.01.2006", search)
	if err != nil {
		tasks, err := in.taskDB.ReadByKeyword(search, 10)
		if err != nil {

			return models.Tasks{}, err
		}

		return tasks, nil
	}

	tasks, err := in.taskDB.ReadByDate(date.Format(nextdate.DateFormat), 10)
	if err != nil {
		return models.Tasks{}, err
	}

	return tasks, nil
}

func (in *TaskService) Create(task *models.Task) (models.TaskCreateResponse, error) {
	if task.Title == "" {
		return models.TaskCreateResponse{}, fmt.Errorf("title empty")
	}

	if task.Date != "" && !dateIsValid(task.Date) {
		return models.TaskCreateResponse{}, fmt.Errorf("fail create task. date invalid")
	}

	if task.Repeat == "" || task.Date == "" || task.Date == time.Now().Format(nextdate.DateFormat) {
		task.Date = time.Now().Format(nextdate.DateFormat)
	} else {
		next, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return models.TaskCreateResponse{}, err
		}
		task.Date = next
	}

	index, err := in.taskDB.Create(task)
	if err != nil {
		return models.TaskCreateResponse{}, err
	}

	return models.TaskCreateResponse{ID: index}, nil
}

func (in *TaskService) Read(id string) (models.Task, error) {
	task, err := in.taskDB.Read(id)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (in *TaskService) Update(task *models.Task) error {
	if task.Title == "" {
		return fmt.Errorf("title empty")
	}

	if task.Date != "" && !dateIsValid(task.Date) {
		return fmt.Errorf("fail create task. date invalid")
	}

	if task.Repeat == "" || task.Date == "" || task.Date == time.Now().Format(nextdate.DateFormat) {
		task.Date = time.Now().Format(nextdate.DateFormat)
	} else {
		_, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	err := in.taskDB.Update(task)
	if err != nil {
		return err
	}

	return nil
}

func (in *TaskService) Delete(id string) error {
	if id == "" {
		return fmt.Errorf("fail delete task. id is empty")
	}

	err := in.taskDB.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (in *TaskService) Done(id string) error {
	task, err := in.Read(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		return in.Delete(id)
	}

	date, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return err
	}

	task.Date = date
	return in.Update(&task)
}

func dateIsValid(date string) bool {
	_, err := time.Parse(nextdate.DateFormat, date)
	return err == nil
}
