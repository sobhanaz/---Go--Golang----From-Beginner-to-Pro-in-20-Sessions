# TaskFlow 📋

A production-style **REST API for managing tasks**, built in Go with the standard
library, JWT authentication, and SQLite. Each user registers, logs in, and
manages their own private list of tasks.

> Built as the capstone project of a 20-session Go course. It demonstrates clean
> architecture, authentication, middleware, testing, and containerized deployment.

---

## ✨ Features

- **RESTful JSON API** built on Go's standard library `net/http` (no web framework)
- **JWT authentication** — register / login, with bcrypt-hashed passwords
- **Per-user data isolation** — every task is scoped to its owner at the SQL layer
- **Middleware** — structured request logging, panic recovery, and auth
- **Clean layered architecture** — `handler → repository → database`, wired via interfaces
- **Configurable** entirely through environment variables
- **Tested** with `httptest` integration tests against a real (temporary) database
- **Dockerized** as a tiny (~15 MB) static image on a distroless base
- **Graceful shutdown** on SIGINT/SIGTERM

## 🧱 Tech stack

| Concern | Choice |
|---------|--------|
| Language | Go 1.25 |
| HTTP | standard library `net/http` (method+path routing, Go 1.22+) |
| Database | SQLite via `modernc.org/sqlite` (pure Go, **no cgo**) |
| Auth | `golang-jwt/jwt/v5` + `golang.org/x/crypto/bcrypt` |
| Tests | standard library `testing` + `net/http/httptest` |
| Deploy | multi-stage Docker → `distroless/static` |

## 🏗️ Architecture

```
main.go                       entry point: config, wire layers, graceful shutdown
internal/
├── config/  config.go        load settings from environment variables
├── models/  task.go,user.go  domain types + request payloads
├── auth/    password.go      bcrypt password hashing
│            jwt.go           issue & verify JWTs
├── store/   store.go         DB connection + schema migration
│            task_store.go    TaskStore repository (user-scoped SQL)
│            user_store.go    UserStore repository
└── api/     server.go        Server, router, repository interfaces, JSON helpers
              middleware.go    Logging, Recovery, Auth middleware
              auth.go          register / login handlers
              tasks.go         task CRUD handlers
              tasks_test.go    integration tests
```

The dependency arrow points one way: `api → store → database`. The `api` package
defines the `TaskRepository` / `UserRepository` interfaces it needs, and the
`store` types satisfy them — so handlers are tested without touching production code.

## 🚀 Run it

### Locally

```bash
cd taskflow
go run .          # http://localhost:8080, creates taskflow.db on first run
```

### With Docker

```bash
docker build -t taskflow .
docker run --rm -p 8080:8080 -e TASKFLOW_JWT_SECRET=change-me taskflow
```

### Configuration

| Env var | Default | Purpose |
|---------|---------|---------|
| `TASKFLOW_ADDR` | `:8080` | listen address |
| `TASKFLOW_DB` | `taskflow.db` | SQLite file path |
| `TASKFLOW_JWT_SECRET` | `dev-secret-change-me` | JWT signing key (**set in prod**) |

## 📡 API

### Auth (public)

| Method | Path | Body | Success |
|--------|------|------|---------|
| POST | `/auth/register` | `{"email","password"}` | 201 + token |
| POST | `/auth/login` | `{"email","password"}` | 200 + token |

### Tasks (require `Authorization: Bearer <token>`)

| Method | Path | Body | Success |
|--------|------|------|---------|
| GET | `/health` | — | 200 (public) |
| GET | `/tasks` | — | 200 |
| POST | `/tasks` | `{"title":"..."}` | 201 |
| GET | `/tasks/{id}` | — | 200 / 404 |
| PUT | `/tasks/{id}` | `{"title":"...","done":true}` | 200 / 404 |
| DELETE | `/tasks/{id}` | — | 204 / 404 |

### Example session

```bash
# 1. Register and capture the token
TOKEN=$(curl -s -X POST localhost:8080/auth/register \
  -d '{"email":"me@example.com","password":"secret123"}' | jq -r .token)

# 2. Create a task
curl -X POST localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" -d '{"title":"Ship TaskFlow"}'

# 3. List your tasks
curl localhost:8080/tasks -H "Authorization: Bearer $TOKEN"
```

## ✅ Test

```bash
go test ./...            # all tests
go test -cover ./...     # with coverage
```

Tests cover registration/login, 401 on missing/invalid tokens, full task CRUD,
and **per-user isolation** (one user cannot read another's tasks).

## 📝 License

MIT — built for learning and portfolio use.
