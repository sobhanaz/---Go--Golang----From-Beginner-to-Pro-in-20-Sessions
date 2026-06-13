> 🌐 **Language / زبان:** English (this file) · [فارسی](session-18.fa.md)

# Session 18 — REST API + Database 🗄️

**Goal (1 hour):** Start building **TaskFlow**, your portfolio project. You'll set
up a proper Go project layout, connect to a real **SQLite database**, and build
full **CRUD** endpoints (Create, Read, Update, Delete) using a clean, layered
architecture. By the end you'll have a working, database-backed REST API.

> **Recap from Session 17:** you can build and test JSON HTTP handlers. Now we
> swap the in-memory map for a real database and organize the code like a pro.

> 📁 **The project lives in [`taskflow/`](../taskflow/)** at the repo root, as its
> own Go module — so you can copy that folder straight to GitHub as a standalone
> project for your CV.

---

## 1. Project layout & the `internal/` convention (10 min)

Real Go projects separate concerns into packages. TaskFlow's layout:

```
taskflow/
├── go.mod                       its own module: "taskflow"
├── main.go                      entry point: open DB, wire layers, serve
└── internal/
    ├── models/   task.go        domain types
    ├── store/    store.go       DB connection + migrations
    │             task_store.go  all SQL for tasks (the "repository")
    └── api/      server.go      router, JSON helpers, repository interface
                  tasks.go       HTTP handlers
                  tasks_test.go  tests
```

> 🔑 **`internal/` is special in Go:** packages under `internal/` can only be
> imported by code in the *same module*. It's the language-enforced way to keep
> your implementation private. Use it for everything that isn't a public library.

**The layered architecture** (the key idea):

```
HTTP request → api (handlers) → store (repository) → database
```

Each layer only talks to the one below it. Handlers don't write SQL; the store
doesn't know about HTTP. This separation makes each piece simple, swappable, and
testable.

---

## 2. Connecting to the database with `database/sql` (15 min)

Go's standard `database/sql` package is a generic interface to SQL databases.
You pair it with a **driver** for your specific database. We use a pure-Go SQLite
driver (no C compiler needed):

```go
import (
    "database/sql"
    _ "modernc.org/sqlite"   // blank import: registers the "sqlite" driver
)

db, err := sql.Open("sqlite", "taskflow.db")
err = db.Ping()              // verify the connection actually works
```

> 🔑 **The blank import `_ "modernc.org/sqlite"`** runs the driver's `init()` to
> register it with `database/sql`, but you don't reference the package directly.
> This driver-registration pattern is standard for SQL in Go.

Add the driver to your module once:

```bash
cd taskflow
go get modernc.org/sqlite
```

### Migrations — creating the schema

On startup we ensure the table exists:

```go
const schema = `
CREATE TABLE IF NOT EXISTS tasks (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    title      TEXT    NOT NULL,
    done       INTEGER NOT NULL DEFAULT 0,
    created_at TEXT    NOT NULL
);`
db.Exec(schema)
```

See [`taskflow/internal/store/store.go`](../taskflow/internal/store/store.go).

---

## 3. The repository: all SQL in one place (20 min)

The **repository pattern** puts every database query for an entity in one type.
The rest of the app calls its methods and never sees SQL.

```go
type TaskStore struct { db *sql.DB }

func (s *TaskStore) Create(title string) (models.Task, error) {
    res, err := s.db.Exec(
        `INSERT INTO tasks (title, done, created_at) VALUES (?, 0, ?)`,
        title, time.Now().UTC().Format(time.RFC3339),
    )
    id, _ := res.LastInsertId()
    return models.Task{ID: id, Title: title, CreatedAt: ...}, nil
}
```

Two `database/sql` essentials:

- **`db.Exec`** — for `INSERT`/`UPDATE`/`DELETE` (no rows returned). Gives you
  `LastInsertId()` and `RowsAffected()`.
- **`db.Query`** (many rows) and **`db.QueryRow`** (one row) — for `SELECT`. You
  read columns with `rows.Scan(&a, &b, ...)`.

> ⚠️ **Always use `?` placeholders** for values — *never* build SQL by string
> concatenation. Placeholders prevent **SQL injection**, the classic security
> hole. `Exec("... WHERE id = ?", id)` is safe; `Exec("... WHERE id = " + id)`
> is dangerous. This is a non-negotiable habit.

When a row isn't found, `QueryRow(...).Scan(...)` returns `sql.ErrNoRows`. We
translate that into our own sentinel `ErrNotFound` (Session 11) so the API layer
can map it to a `404`:

```go
if errors.Is(err, sql.ErrNoRows) {
    return models.Task{}, ErrNotFound
}
```

Study the full repository: [`taskflow/internal/store/task_store.go`](../taskflow/internal/store/task_store.go).

---

## 4. Wiring handlers to the repository via an interface (15 min)

The API layer defines the behavior it needs as an **interface**, then depends on
that — not on the concrete `TaskStore`. This is dependency inversion, and it's
why the handlers are testable (Session 10 in action):

```go
// in package api
type TaskRepository interface {
    Create(title string) (models.Task, error)
    List() ([]models.Task, error)
    Get(id int64) (models.Task, error)
    Update(id int64, title string, done bool) (models.Task, error)
    Delete(id int64) error
}

type Server struct { tasks TaskRepository }
```

`*store.TaskStore` satisfies this interface automatically (implicit
satisfaction!). Handlers map repository results and errors to HTTP:

```go
task, err := s.tasks.Get(id)
if errors.Is(err, store.ErrNotFound) {
    writeError(w, http.StatusNotFound, "task not found")
    return
}
```

`main.go` wires the three layers together:

```go
db, _ := store.Open("taskflow.db")
taskStore := store.NewTaskStore(db)        // store layer
server := api.NewServer(taskStore)         // api layer (gets the repo)
http.ListenAndServe(":8080", server.Routes())
```

### Run and test it

```bash
cd taskflow
go run .                       # real server + real DB
go test ./...                  # integration tests against a temp DB
```

The tests use a **real SQLite DB in `t.TempDir()`** — proving the whole stack
(handler → SQL → database) works, then cleaning up automatically. See
[`taskflow/internal/api/tasks_test.go`](../taskflow/internal/api/tasks_test.go).

> 💡 **Why a real DB in tests, not a fake?** For a small project, exercising real
> SQL catches more bugs and is simple with SQLite's temp/in-memory support. The
> `TaskRepository` interface still lets you swap a fake in if you ever need to.

---

## 🎯 Exercises (do these before Session 19!)

Work inside `taskflow/`:

1. **Run the whole flow:** Start the server and use `curl` to create, list,
   update, get, and delete a task. Watch `taskflow.db` appear.
2. **Add a field:** Add a `Priority int` column (migration + model + SQL in
   Create/Update/Scan). Confirm it round-trips through the API.
3. **Filter endpoint:** Add `GET /tasks?done=true` that lists only completed
   tasks (read the query param, add a `WHERE done = ?` query method).
4. **Validation:** Reject titles longer than 200 characters with a `400`.
5. **More tests:** Add a test for the update-not-found case (`PUT /tasks/999`
   → 404) and for your new filter endpoint.

---

## ✅ Session 18 Checklist

- [ ] I understand the layered layout (api → store → db) and `internal/`
- [ ] I can open a SQLite DB with `database/sql` + a driver blank-import
- [ ] I run a migration to create the schema on startup
- [ ] I can write `Exec`, `Query`, and `QueryRow` with `?` placeholders
- [ ] I know why placeholders prevent SQL injection
- [ ] I translate `sql.ErrNoRows` into my own `ErrNotFound`
- [ ] I depend on a `TaskRepository` interface in the API layer
- [ ] I can run the server and the integration tests
- [ ] I completed all 5 exercises

**Previous:** [← Session 17](session-17.md) · **Next:** [Session 19 — Auth, Middleware & Config →](session-19.md)
