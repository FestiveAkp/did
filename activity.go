package main

import (
	"database/sql"
	"time"
)

type Activity struct {
	ID        int64
	TaskID    int64
	Note      string
	CreatedAt time.Time
}

type ActivityStore struct {
	db *sql.DB
}

func NewActivityStore(db *sql.DB) *ActivityStore {
	return &ActivityStore{db: db}
}

func (s *ActivityStore) Create(taskID int64, note string) (Activity, error) {
	res, err := s.db.Exec(
		`INSERT INTO activities (task_id, note) VALUES (?, ?)`,
		taskID, note,
	)
	if err != nil {
		return Activity{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return Activity{}, err
	}
	row := s.db.QueryRow(
		`SELECT id, task_id, note, created_at FROM activities WHERE id = ?`, id,
	)
	return scanActivity(row)
}

func (s *ActivityStore) ListForTask(taskID int64) ([]Activity, error) {
	rows, err := s.db.Query(
		`SELECT id, task_id, note, created_at FROM activities WHERE task_id = ? ORDER BY created_at ASC`,
		taskID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []Activity
	for rows.Next() {
		a, err := scanActivity(rows)
		if err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}
	return activities, rows.Err()
}

func (s *ActivityStore) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM activities WHERE id = ?`, id)
	return err
}

func scanActivity(row interface{ Scan(dest ...any) error }) (Activity, error) {
	var a Activity
	if err := row.Scan(&a.ID, &a.TaskID, &a.Note, &a.CreatedAt); err != nil {
		return Activity{}, err
	}
	return a, nil
}
