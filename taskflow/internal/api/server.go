// Package api wires HTTP handlers to the repositories.
package api

import (
	"encoding/json"
	"net/http"

	"taskflow/internal/models"
)

// TaskRepository is the behavior the API needs from task storage.
// Every method is scoped to a userID so users only see their own tasks.
type TaskRepository interface {
	Create(userID int64, title string) (models.Task, error)
	List(userID int64) ([]models.Task, error)
	Get(userID, id int64) (models.Task, error)
	Update(userID, id int64, title string, done bool) (models.Task, error)
	Delete(userID, id int64) error
}

// UserRepository is the behavior the API needs from user storage.
type UserRepository interface {
	Create(email, passwordHash string) (models.User, error)
	GetByEmail(email string) (models.User, error)
}

// Server holds the API's dependencies.
type Server struct {
	tasks     TaskRepository
	users     UserRepository
	jwtSecret string
}

func NewServer(tasks TaskRepository, users UserRepository, jwtSecret string) *Server {
	return &Server{tasks: tasks, users: users, jwtSecret: jwtSecret}
}

// Routes builds the router and applies middleware.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	// Public routes (no auth required).
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("POST /auth/register", s.handleRegister)
	mux.HandleFunc("POST /auth/login", s.handleLogin)

	// Protected task routes: wrap each with the Auth middleware.
	mux.Handle("GET /tasks", s.Auth(http.HandlerFunc(s.handleListTasks)))
	mux.Handle("POST /tasks", s.Auth(http.HandlerFunc(s.handleCreateTask)))
	mux.Handle("GET /tasks/{id}", s.Auth(http.HandlerFunc(s.handleGetTask)))
	mux.Handle("PUT /tasks/{id}", s.Auth(http.HandlerFunc(s.handleUpdateTask)))
	mux.Handle("DELETE /tasks/{id}", s.Auth(http.HandlerFunc(s.handleDeleteTask)))

	// Global middleware wraps the whole mux. Recovery is outermost so it
	// catches panics from everything, including Logging.
	return Recovery(Logging(mux))
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// --- small response helpers shared by all handlers ---

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
