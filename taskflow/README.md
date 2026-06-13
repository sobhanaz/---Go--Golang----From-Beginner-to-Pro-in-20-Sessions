# TaskFlow 📋

A small, production-style **REST API for managing tasks**, built in Go with the
standard library and SQLite. This is the portfolio project for the
[Go course](../README.md) — built across Sessions 18–20.

> **Status:** Session 18 — task CRUD with database persistence.
> Sessions 19–20 add authentication, middleware, config, Docker, and polish.

## Features (so far)

- RESTful JSON API for tasks (create, read, update, delete)
- SQLite persistence via a pure-Go driver (no cgo)
- Clean layered architecture: **handler → repository → database**
- Handlers depend on an interface, making them easy to test
- Integration tests using `httptest` against a real (temporary) database

## Architecture

```
main.go                      entry point: open DB, wire layers, start server
internal/
├── models/   task.go        domain types (Task) + input payloads
├── store/    store.go       DB connection + schema migration
│             task_store.go  TaskStore repository (all SQL lives here)
└── api/      server.go      Server, router, JSON helpers, TaskRepository interface
              tasks.go       HTTP handlers for /tasks
              tasks_test.go  handler/integration tests
```

The dependency arrow points one way: `api` → `store` → database. The `api`
package defines the `TaskRepository` interface it needs, and `*store.TaskStore`
satisfies it. This is why tests can run without changing any production code.

## Run it

```bash
cd taskflow
go run .                     # starts on http://localhost:8080, creates taskflow.db
```

## API

| Method | Path          | Body                          | Success |
|--------|---------------|-------------------------------|---------|
| GET    | `/health`     | —                             | 200     |
| GET    | `/tasks`      | —                             | 200     |
| POST   | `/tasks`      | `{"title":"..."}`             | 201     |
| GET    | `/tasks/{id}` | —                             | 200/404 |
| PUT    | `/tasks/{id}` | `{"title":"...","done":true}` | 200/404 |
| DELETE | `/tasks/{id}` | —                             | 204/404 |

### Examples

```bash
curl localhost:8080/health
curl -X POST localhost:8080/tasks -d '{"title":"Learn Go"}'
curl localhost:8080/tasks
curl -X PUT localhost:8080/tasks/1 -d '{"title":"Learn Go","done":true}'
curl -X DELETE localhost:8080/tasks/1
```

## Test

```bash
cd taskflow
go test ./...
go test -cover ./...
```

## Tech

- Go (standard library `net/http`, `database/sql`, `encoding/json`)
- SQLite via `modernc.org/sqlite` (pure Go, no cgo)
