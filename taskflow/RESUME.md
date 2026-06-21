# TaskFlow — Résumé & Interview Notes

Copy-paste these into your CV / LinkedIn, then tailor the wording to the role.

---

## Résumé bullets (concise)

**TaskFlow — Task Management REST API (Go)** · github.com/sobhanaz/go-from-beginner-to-pro/tree/main/taskflow

- Built a multi-user task-management **REST API in Go** using the standard
  library `net/http`, with **JWT authentication** and bcrypt-hashed passwords.
- Designed a **clean, layered architecture** (handler → repository → database)
  decoupled through interfaces, enabling fast handler tests against a real DB.
- Implemented **middleware** for structured request logging, panic recovery, and
  auth, plus **per-user data isolation** enforced at the SQL layer (`WHERE user_id = ?`).
- Added a **task-priority** feature (low/medium/high) with server-side validation
  and **composable query filters** (`?done=`, `?priority=`), wired end-to-end from
  model and migration through repository, handler, and tests.
- Wrote **integration tests** with `net/http/httptest` (auth, CRUD, access
  control, validation, filters); 10 passing tests against a temporary database.
- **Containerized** with a multi-stage Docker build producing a ~15 MB static
  **distroless** image, with **graceful shutdown** (SIGTERM + `context` timeout)
  for zero-downtime deploys.

## One-line version (for a "Projects" header)

> *TaskFlow* — a JWT-authenticated, multi-user task REST API in Go (net/http,
> SQLite), with middleware, per-user isolation, prioritized tasks with filters,
> integration tests, and a Dockerized distroless deploy.

---

## Interview talking points (know the *why*)

- **Why no web framework?** Go's standard library `net/http` (with 1.22+
  method/path routing) is enough for a clean REST API — fewer dependencies, less
  magic, easier to reason about.
- **How does authentication work?** Stateless **JWTs**: the signed token carries
  the user's ID in its subject, so the server verifies the signature instead of
  storing server-side sessions. Passwords are bcrypt-hashed, never stored or
  logged in plaintext.
- **How is one user's data kept private?** Every task query is scoped with
  `WHERE user_id = ?`. A user physically cannot read another's rows — proven by a
  test where one user gets `404` for another's task ID.
- **Walk me through adding the priority feature.** Touched every layer: a
  `Priority` field + validation in the model, an idempotent `ALTER TABLE`
  migration, repository methods and a dynamic-but-parameterized filter query,
  handler validation returning `400` on bad input, and table-style tests.
- **How do you prevent SQL injection?** All values are **bound parameters** (`?`),
  never string-concatenated — even in the dynamic filter query, which appends
  `AND` clauses but passes the values as args.
- **How would you scale it?** Swap SQLite for PostgreSQL behind the same
  `TaskRepository` interface; the handlers and tests don't change. That decoupling
  is the payoff of defining the interface in the consumer (`api`) package.
- **How do you test handlers?** `httptest` fires requests at the router with no
  real network, against a fresh temporary database per test — fast, isolated,
  deterministic.
- **Why graceful shutdown?** Orchestrators (e.g. Kubernetes) send `SIGTERM`
  before stopping a container; the server stops accepting new requests, lets
  in-flight ones finish within a `context` deadline, then exits cleanly.
