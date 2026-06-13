package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"taskflow/internal/auth"
	"taskflow/internal/models"
	"taskflow/internal/store"
)

// POST /auth/register — create an account and return a token.
func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var in models.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	if in.Email == "" || len(in.Password) < 6 {
		writeError(w, http.StatusBadRequest, "email required and password must be at least 6 chars")
		return
	}

	hash, err := auth.HashPassword(in.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	user, err := s.users.Create(in.Email, hash)
	if errors.Is(err, store.ErrEmailTaken) {
		writeError(w, http.StatusConflict, "email already registered")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	s.respondWithToken(w, http.StatusCreated, user)
}

// POST /auth/login — verify credentials and return a token.
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var in models.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	user, err := s.users.GetByEmail(in.Email)
	// Use the SAME generic message for "no user" and "wrong password" so we
	// don't reveal which emails are registered (a small security best practice).
	if err != nil || !auth.CheckPassword(user.PasswordHash, in.Password) {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	s.respondWithToken(w, http.StatusOK, user)
}

// respondWithToken issues a JWT for the user and writes it as JSON.
func (s *Server) respondWithToken(w http.ResponseWriter, status int, user models.User) {
	token, err := auth.GenerateToken(s.jwtSecret, user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create token")
		return
	}
	writeJSON(w, status, map[string]any{
		"token": token,
		"user":  user, // PasswordHash is hidden by its json:"-" tag
	})
}
