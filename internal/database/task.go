package database

import (
	"database/sql"
	"fmt"

	"github.com/Vafeda/TODO-List/internal/models"
)

type TaskDB struct {
	db *sql.DB
}

func NewTaskDB(db *sql.DB) *TaskDB {
	return &TaskDB{
		db: db,
	}
}

func (t *TaskDB) ReadAll(limit int) (models.Tasks, error) {
	rows, err := t.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY DATE LIMIT $1", limit)
	if err != nil {
		return models.Tasks{}, err
	}
	defer rows.Close()

	tasks := make([]models.Task, 0, limit)
	for rows.Next() {
		task := models.Task{}
		err = rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return models.Tasks{}, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return models.Tasks{}, err
	}

	return models.Tasks{Tasks: tasks}, nil
}

func (t *TaskDB) ReadByDate(date string, limit int) (models.Tasks, error) {
	rows, err := t.db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = $1 LIMIT $1", date, limit)
	if err != nil {
		return models.Tasks{}, err
	}
	defer rows.Close()

	tasks := make([]models.Task, 0, limit)
	for rows.Next() {
		task := models.Task{}
		err = rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return models.Tasks{}, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return models.Tasks{}, err
	}

	return models.Tasks{Tasks: tasks}, nil
}

func (t *TaskDB) ReadByKeyword(key string, limit int) (models.Tasks, error) {
	rows, err := t.db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE $1 OR comment LIKE $1 ORDER BY date LIMIT $2", key, limit)
	if err != nil {
		return models.Tasks{}, err
	}
	defer rows.Close()

	tasks := make([]models.Task, 0, limit)
	for rows.Next() {
		task := models.Task{}
		err = rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return models.Tasks{}, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return models.Tasks{}, err
	}

	return models.Tasks{Tasks: tasks}, nil
}

func (t *TaskDB) Create(task *models.Task) (int64, error) {
	r, err := t.db.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4)`,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat)

	if err != nil {
		return 0, err
	}

	index, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return index, nil
}

func (t *TaskDB) Read(id string) (models.Task, error) {
	var task models.Task

	err := t.db.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = $1`,
		id,
	).Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Task{}, fmt.Errorf("задача не найдена")
		}
		return models.Task{}, fmt.Errorf("ошибка базы данных: %w", err)
	}

	return task, nil
}

func (t *TaskDB) Update(task *models.Task) error {
	r, err := t.db.Exec(`UPDATE scheduler SET date = $1, title = $2, comment = $3, repeat = $4 WHERE id = $5`,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		task.Id)

	if err != nil {
		return err
	}

	count, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	return nil
}

func (t *TaskDB) Delete(id string) error {
	r, err := t.db.Exec(`DELETE FROM scheduler WHERE id = $1`, id)
	if err != nil {
		return err
	}

	count, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for deleting task`)
	}

	return nil
}
