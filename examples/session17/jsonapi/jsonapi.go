// Session 17 — a small JSON REST API with in-memory storage.
// This previews the structure of the final TaskFlow project.
// Run:  go run examples/session17/jsonapi/jsonapi.go
// Try:
//   curl localhost:8080/health
//   curl localhost:8080/tasks
//   curl -X POST localhost:8080/tasks -d '{"title":"Learn Go"}'
//   curl localhost:8080/tasks/1
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Task is our domain model. JSON tags control the wire format.
type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// Server holds dependencies (here, an in-memory store guarded by a mutex).
// In the final project this is where a database connection would live.
type Server struct {
	mu     sync.Mutex
	tasks  map[int]Task
	nextID int
}

func NewServer() *Server {
	return &Server{tasks: make(map[int]Task), nextID: 1}
}

// writeJSON is a small helper: set the header, status, and encode the body.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// GET /health — a liveness check.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// GET /tasks — list all tasks.
func (s *Server) handleListTasks(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	list := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		list = append(list, t)
	}
	writeJSON(w, http.StatusOK, list)
}

// POST /tasks — create a task from a JSON body.
func (s *Server) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if input.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}

	s.mu.Lock()
	task := Task{ID: s.nextID, Title: input.Title}
	s.tasks[task.ID] = task
	s.nextID++
	s.mu.Unlock()

	writeJSON(w, http.StatusCreated, task) // 201 Created
}

// GET /tasks/{id} — fetch one task by path parameter.
func (s *Server) handleGetTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id")) // Go 1.22+ path wildcard
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	s.mu.Lock()
	task, ok := s.tasks[id]
	s.mu.Unlock()
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}
	writeJSON(w, http.StatusOK, task)
}

// routes wires patterns to handlers and returns the mux (an http.Handler).
// Keeping this separate makes it easy to test (see jsonapi_test.go).
func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /tasks", s.handleListTasks)
	mux.HandleFunc("POST /tasks", s.handleCreateTask)
	mux.HandleFunc("GET /tasks/{id}", s.handleGetTask)
	return mux
}

func main() {
	srv := NewServer()
	addr := ":8080"
	log.Printf("TaskFlow preview API listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, srv.routes()); err != nil {
		log.Fatal(err)
	}
}
