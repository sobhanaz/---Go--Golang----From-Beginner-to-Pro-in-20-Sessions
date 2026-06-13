// TaskFlow — a small, authenticated REST API for managing tasks.
// Run:  go run .            (from the taskflow/ directory)
// Config via env: TASKFLOW_ADDR, TASKFLOW_DB, TASKFLOW_JWT_SECRET
package main

import (
	"log"
	"net/http"

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

	log.Printf("TaskFlow listening on http://localhost%s", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, server.Routes()); err != nil {
		log.Fatal(err)
	}
}
