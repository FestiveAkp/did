package main

import (
	"database/sql"
	"time"
)

type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

// AllStatuses lists the statuses a task can be set to, in display order.
var AllStatuses = []Status{StatusTodo, StatusInProgress, StatusDone}

func (s Status) Label() string {
	switch s {
	case StatusInProgress:
		return "In Progress"
	case StatusDone:
		return "Done"
	default:
		return "To Do"
	}
}

// Icon returns a circle glyph representing the status.
func (s Status) Icon() string {
	switch s {
	case StatusInProgress:
		return "◐"
	case StatusDone:
		return "●"
	default:
		return "○"
	}
}

type Task struct {
	ID        int64
	Title     string
	Status    Status
	CreatedAt time.Time
}

type TaskStore struct {
	db *sql.DB
}

func NewTaskStore(db *sql.DB) *TaskStore {
	return &TaskStore{db: db}
}

func (s *TaskStore) Create(title string) (Task, error) {
	res, err := s.db.Exec(
		`INSERT INTO tasks (title, status) VALUES (?, ?)`,
		title, StatusTodo,
	)
	if err != nil {
		return Task{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return Task{}, err
	}
	return s.Get(id)
}

func (s *TaskStore) Get(id int64) (Task, error) {
	row := s.db.QueryRow(
		`SELECT id, title, status, created_at FROM tasks WHERE id = ?`, id,
	)
	return scanTask(row)
}

func (s *TaskStore) List() ([]Task, error) {
	rows, err := s.db.Query(
		`SELECT id, title, status, created_at FROM tasks ORDER BY created_at ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (s *TaskStore) SetStatus(id int64, status Status) error {
	_, err := s.db.Exec(`UPDATE tasks SET status = ? WHERE id = ?`, status, id)
	return err
}

func (s *TaskStore) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	return err
}

func scanTask(row interface{ Scan(dest ...any) error }) (Task, error) {
	var t Task
	var status string
	if err := row.Scan(&t.ID, &t.Title, &status, &t.CreatedAt); err != nil {
		return Task{}, err
	}
	t.Status = Status(status)
	return t, nil
}
