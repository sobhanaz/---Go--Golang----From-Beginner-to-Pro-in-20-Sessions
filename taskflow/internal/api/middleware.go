package api

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"taskflow/internal/auth"
)

// A middleware wraps a handler to add behavior, returning a new handler.
// The signature func(http.Handler) http.Handler is the Go middleware idiom.

// ctxKey is an unexported type for context keys, so other packages can't collide.
type ctxKey string

const userIDKey ctxKey = "userID"

// Logging records the method, path, status, and duration of each request.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Wrap the ResponseWriter to capture the status code.
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, rec.status, time.Since(start))
	})
}

// Recovery catches a panic in any handler so one bad request can't crash the
// whole server (Session 11's recover, applied as middleware).
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				writeError(w, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Auth verifies the Bearer token and stores the user ID in the request context.
// Requests without a valid token are rejected with 401.
func (s *Server) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(header, prefix) {
			writeError(w, http.StatusUnauthorized, "missing or malformed Authorization header")
			return
		}
		token := strings.TrimPrefix(header, prefix)
		userID, err := auth.ParseToken(s.jwtSecret, token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		// Attach the user ID to the context for handlers to read.
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// userIDFromContext retrieves the authenticated user's ID set by Auth.
func userIDFromContext(r *http.Request) int64 {
	id, _ := r.Context().Value(userIDKey).(int64)
	return id
}

// statusRecorder remembers the status code written to the response.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
