// Package config loads application settings from environment variables,
// falling back to sensible defaults. This is how real apps stay configurable
// across dev/staging/production without code changes (Session 14's os.Getenv).
package config

import "os"

type Config struct {
	Addr        string // e.g. ":8080"
	DatabaseDSN string // e.g. "taskflow.db"
	JWTSecret   string // signing key for JWTs — MUST be set in production
}

// Load reads config from the environment with defaults.
func Load() Config {
	return Config{
		Addr:        getenv("TASKFLOW_ADDR", ":8080"),
		DatabaseDSN: getenv("TASKFLOW_DB", "taskflow.db"),
		// The default secret is only for local dev. Always override it in prod.
		JWTSecret: getenv("TASKFLOW_JWT_SECRET", "dev-secret-change-me"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
