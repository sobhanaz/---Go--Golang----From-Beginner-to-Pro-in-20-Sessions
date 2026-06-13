package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"taskflow/internal/models"
	"taskflow/internal/store"
)

// parseID reads the {id} path wildcard and converts it to int64.
func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}

// GET /tasks
func (s *Server) handleListTasks(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	tasks, err := s.tasks.List(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list tasks")
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

// POST /tasks
func (s *Server) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	var in models.CreateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if strings.TrimSpace(in.Title) == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	task, err := s.tasks.Create(userID, in.Title)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create task")
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

// GET /tasks/{id}
func (s *Server) handleGetTask(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	task, err := s.tasks.Get(userID, id)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not get task")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

// PUT /tasks/{id}
func (s *Server) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in models.UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if strings.TrimSpace(in.Title) == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	task, err := s.tasks.Update(userID, id, in.Title, in.Done)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not update task")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

// DELETE /tasks/{id}
func (s *Server) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	err = s.tasks.Delete(userID, id)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not delete task")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
