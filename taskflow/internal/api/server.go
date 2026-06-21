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
	Create(userID int64, title, priority string) (models.Task, error)
	List(userID int64, filter models.TaskFilter) ([]models.Task, error)
	Get(userID, id int64) (models.Task, error)
	Update(userID, id int64, title, priority string, done bool) (models.Task, error)
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
	// "GET /{$}" matches ONLY the exact root path "/" (Go 1.22+), so it acts as
	// a friendly index without becoming a catch-all for unknown paths.
	mux.HandleFunc("GET /{$}", s.handleIndex)
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

// handleIndex returns a small JSON description of the API at the root path,
// so opening "/" in a browser shows what's available instead of a 404.
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"service": "TaskFlow API",
		"version": "1.0",
		"endpoints": map[string]string{
			"GET  /health":        "liveness check (public)",
			"POST /auth/register": "create an account, returns a JWT (public)",
			"POST /auth/login":    "log in, returns a JWT (public)",
			"GET  /tasks":         "list your tasks; filters: ?done=, ?priority= (auth)",
			"POST /tasks":         "create a task (auth)",
			"GET  /tasks/{id}":    "get one task (auth)",
			"PUT  /tasks/{id}":    "update a task (auth)",
			"DELETE /tasks/{id}":  "delete a task (auth)",
		},
		"docs": "Send 'Authorization: Bearer <token>' for /tasks endpoints.",
	})
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
