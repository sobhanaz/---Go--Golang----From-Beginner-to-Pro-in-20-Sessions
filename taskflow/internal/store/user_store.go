package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"taskflow/internal/models"
)

// ErrUserNotFound is returned when no user matches a lookup.
var ErrUserNotFound = errors.New("user not found")

// ErrEmailTaken is returned when registering with an already-used email.
var ErrEmailTaken = errors.New("email already registered")

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

// Create inserts a new user. The caller passes an already-hashed password.
func (s *UserStore) Create(email, passwordHash string) (models.User, error) {
	now := time.Now().UTC()
	res, err := s.db.Exec(
		`INSERT INTO users (email, password_hash, created_at) VALUES (?, ?, ?)`,
		email, passwordHash, now.Format(time.RFC3339),
	)
	if err != nil {
		// The UNIQUE constraint on email fails here for duplicates.
		return models.User{}, ErrEmailTaken
	}
	id, _ := res.LastInsertId()
	return models.User{ID: id, Email: email, PasswordHash: passwordHash, CreatedAt: now}, nil
}

// GetByEmail looks up a user by email (used during login).
func (s *UserStore) GetByEmail(email string) (models.User, error) {
	row := s.db.QueryRow(
		`SELECT id, email, password_hash, created_at FROM users WHERE email = ?`, email)
	var (
		u         models.User
		createdAt string
	)
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &createdAt)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		return models.User{}, fmt.Errorf("get user: %w", err)
	}
	u.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return u, nil
}
