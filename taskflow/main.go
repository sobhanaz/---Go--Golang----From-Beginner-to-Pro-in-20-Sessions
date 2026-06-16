// TaskFlow — a small, authenticated REST API for managing tasks.
// Run:  go run .            (from the taskflow/ directory)
// Config via env: TASKFLOW_ADDR, TASKFLOW_DB, TASKFLOW_JWT_SECRET
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"taskflow/internal/api"
	"taskflow/internal/config"
	"taskflow/internal/store"
)

func main() {
	cfg := config.Load()

	// Open the database (creates the file + schema on first run).
	db, err := store.Open(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	// Wire the layers: stores -> API server (with the JWT secret).
	taskStore := store.NewTaskStore(db)
	userStore := store.NewUserStore(db)
	server := api.NewServer(taskStore, userStore, cfg.JWTSecret)

	// Use an explicit http.Server so we can shut it down gracefully.
	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      server.Routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run the server in a goroutine so main can wait for a shutdown signal.
	go func() {
		log.Printf("TaskFlow listening on http://localhost%s", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Block until we receive an interrupt (Ctrl+C) or termination signal.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	// Give in-flight requests up to 10s to finish, then stop.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
	log.Println("server stopped")
}
