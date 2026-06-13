package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"taskflow/internal/models"
)

// ErrNotFound is returned when a task does not exist (a sentinel error, Session 11).
var ErrNotFound = errors.New("task not found")

// TaskStore is the repository: it owns all SQL for tasks.
// Every method is scoped to a userID so users only ever touch their own tasks.
type TaskStore struct {
	db *sql.DB
}

func NewTaskStore(db *sql.DB) *TaskStore {
	return &TaskStore{db: db}
}

// Create inserts a new task owned by userID.
func (s *TaskStore) Create(userID int64, title string) (models.Task, error) {
	now := time.Now().UTC()
	res, err := s.db.Exec(
		`INSERT INTO tasks (user_id, title, done, created_at) VALUES (?, ?, 0, ?)`,
		userID, title, now.Format(time.RFC3339),
	)
	if err != nil {
		return models.Task{}, fmt.Errorf("insert task: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.Task{}, fmt.Errorf("last insert id: %w", err)
	}
	return models.Task{ID: id, UserID: userID, Title: title, Done: false, CreatedAt: now}, nil
}

// List returns all of userID's tasks, newest first.
func (s *TaskStore) List(userID int64) ([]models.Task, error) {
	rows, err := s.db.Query(
		`SELECT id, user_id, title, done, created_at FROM tasks
		 WHERE user_id = ? ORDER BY id DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

// Get returns one of userID's tasks by ID, or ErrNotFound.
func (s *TaskStore) Get(userID, id int64) (models.Task, error) {
	row := s.db.QueryRow(
		`SELECT id, user_id, title, done, created_at FROM tasks
		 WHERE id = ? AND user_id = ?`, id, userID)
	t, err := scanTask(row)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Task{}, ErrNotFound
	}
	if err != nil {
		return models.Task{}, fmt.Errorf("get task: %w", err)
	}
	return t, nil
}

// Update modifies one of userID's tasks.
func (s *TaskStore) Update(userID, id int64, title string, done bool) (models.Task, error) {
	res, err := s.db.Exec(
		`UPDATE tasks SET title = ?, done = ? WHERE id = ? AND user_id = ?`,
		title, done, id, userID)
	if err != nil {
		return models.Task{}, fmt.Errorf("update task: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return models.Task{}, ErrNotFound
	}
	return s.Get(userID, id)
}

// Delete removes one of userID's tasks.
func (s *TaskStore) Delete(userID, id int64) error {
	res, err := s.db.Exec(`DELETE FROM tasks WHERE id = ? AND user_id = ?`, id, userID)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// scanner is satisfied by both *sql.Row and *sql.Rows, so one helper covers both.
type scanner interface {
	Scan(dest ...any) error
}

func scanTask(sc scanner) (models.Task, error) {
	var (
		t         models.Task
		createdAt string
	)
	if err := sc.Scan(&t.ID, &t.UserID, &t.Title, &t.Done, &createdAt); err != nil {
		return models.Task{}, err
	}
	t.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return t, nil
}
