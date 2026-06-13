// Package store handles the database connection and persistence.
package store

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // pure-Go SQLite driver, registered as "sqlite"
)

// Open opens (or creates) the SQLite database at dsn and runs migrations.
// dsn examples:  "taskflow.db"  or  ":memory:"  for an ephemeral DB.
func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	// Ping verifies the connection is actually usable.
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}

// migrate creates the schema if it does not exist. In larger projects this
// is replaced by versioned migration files; here a single statement is fine.
func migrate(db *sql.DB) error {
	const schema = `
	CREATE TABLE IF NOT EXISTS users (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		email         TEXT    NOT NULL UNIQUE,
		password_hash TEXT    NOT NULL,
		created_at    TEXT    NOT NULL
	);

	CREATE TABLE IF NOT EXISTS tasks (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id    INTEGER NOT NULL,
		title      TEXT    NOT NULL,
		done       INTEGER NOT NULL DEFAULT 0,
		created_at TEXT    NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`
	_, err := db.Exec(schema)
	return err
}
