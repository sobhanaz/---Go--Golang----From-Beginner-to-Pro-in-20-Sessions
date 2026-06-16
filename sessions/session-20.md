> 🌐 **Language / زبان:** English (this file) · [فارسی](session-20.fa.md)

# Session 20 — Polish & Ship 🚢

**Goal (1 hour):** Take TaskFlow from "works on my machine" to "a recruiter can
clone, run, and be impressed." You'll add **graceful shutdown**, **Docker**, a
**Makefile**, and a polished **README** — then turn the whole project into
**CV bullet points** and interview talking points. This is the session that makes
your work *presentable*.

> **Recap from Session 19:** TaskFlow is a secure, multi-user, tested REST API.
> Now we make it deployable and presentable — the difference between a hobby
> script and a portfolio piece.

> 📂 All code is in [`taskflow/`](../taskflow/).

---

## 1. Graceful shutdown (15 min)

A naive server (`http.ListenAndServe(...)`) dies instantly when killed, dropping
any in-flight requests. Production servers **shut down gracefully**: stop
accepting new requests, let current ones finish, then exit. See
[`taskflow/main.go`](../taskflow/main.go):

```go
srv := &http.Server{Addr: cfg.Addr, Handler: server.Routes(), /* timeouts */}

// Run the server in a goroutine so main can wait for a signal.
go func() {
    if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
        log.Fatalf("server error: %v", err)
    }
}()

// Block until an interrupt/termination signal arrives.
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

// Give in-flight requests up to 10s to finish.
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
srv.Shutdown(ctx)
```

This combines everything you learned: **goroutines** (Session 12), **channels**
to wait for a signal, and **context** with a timeout (Session 13). The
`http.Server` also sets `ReadTimeout`/`WriteTimeout`/`IdleTimeout` — basic
hardening so slow clients can't tie up resources.

> 🔑 **Why it matters:** orchestrators like Kubernetes send `SIGTERM` before
> stopping a container. A graceful server finishes its work and exits cleanly
> instead of erroring out mid-request. Interviewers love seeing this.

---

## 2. Dockerizing — a tiny, secure image (20 min)

A **multi-stage** Dockerfile builds in one image and ships in a minimal one. See
[`taskflow/Dockerfile`](../taskflow/Dockerfile):

```dockerfile
# ---- Build stage ----
FROM golang:1.25 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download              # cached unless deps change
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /taskflow .

# ---- Runtime stage ----
FROM gcr.io/distroless/static-debian12
COPY --from=build /taskflow /app/taskflow
EXPOSE 8080
ENTRYPOINT ["/app/taskflow"]
```

Three ideas that make this *good*, not just working:

1. **Multi-stage build.** The heavy Go toolchain stays in the build stage; the
   final image contains only your binary. Result: a ~15 MB image instead of ~1 GB.
2. **`CGO_ENABLED=0`** produces a fully static binary. This is only possible
   because TaskFlow uses the **pure-Go** SQLite driver (`modernc.org/sqlite`) —
   a deliberate choice from Session 18 paying off now.
3. **Distroless base.** `distroless/static` has no shell and no package manager,
   so the attack surface is tiny. `-ldflags="-s -w"` strips debug info to shrink
   the binary further.

```bash
docker build -t taskflow .
docker run --rm -p 8080:8080 -e TASKFLOW_JWT_SECRET=change-me taskflow
```

A [`.dockerignore`](../taskflow/.dockerignore) keeps the build context small, and
a [`.gitignore`](../taskflow/.gitignore) keeps the binary and `*.db` out of git.

---

## 3. Developer ergonomics: the Makefile (10 min)

A [`Makefile`](../taskflow/Makefile) gives anyone a one-word entry to common tasks
— a small touch that signals care:

```bash
make help          # list commands
make run           # go run .
make test          # go test ./...
make cover         # go test -cover ./...
make docker-build  # build the image
```

> 💡 You don't *need* Make, but `make test` / `make run` is friendlier than
> remembering flags, and recruiters see that you think about the next developer.

---

## 4. The README that gets you the interview (10 min)

For a portfolio project, **the README is the most important file** — it's the
first (often only) thing a recruiter reads. A strong one has, in order:

1. **One-line description** of what it is.
2. **Feature bullets** — what it does, with the buzzwords that matter (JWT, REST,
   Docker, tests).
3. **Tech stack** table.
4. **Architecture** — a directory tree + one paragraph on the design.
5. **Run instructions** — local *and* Docker, copy-pasteable.
6. **API reference** — endpoints, with a worked `curl` example.
7. **How to test.**

[`taskflow/README.md`](../taskflow/README.md) is written exactly this way. Study
it as a template; a clear README often impresses more than the code itself,
because most reviewers skim.

---

## 5. Putting TaskFlow on your CV (5 min)

Describe the project with **impact and specifics**, not just "made an API." Pick a
few bullets like these for your résumé:

> **TaskFlow — REST API (Go)** · [github.com/you/taskflow]
> - Built a multi-user task-management REST API in **Go** using the standard
>   library `net/http`, with **JWT authentication** and bcrypt password hashing.
> - Designed a **clean, layered architecture** (handler → repository → database)
>   decoupled via interfaces, enabling handler tests against a real database.
> - Implemented **middleware** for structured logging, panic recovery, and auth,
>   and **per-user data isolation** enforced at the SQL layer.
> - Wrote **integration tests** with `net/http/httptest` covering auth, CRUD, and
>   access-control; achieved meaningful coverage across the API.
> - **Containerized** with a multi-stage Docker build producing a ~15 MB static
>   distroless image, with graceful shutdown for zero-downtime deploys.

### Interview talking points (be ready to explain *why*)

- **"Why no framework?"** — Go's standard library is enough for a clean REST API;
  fewer dependencies, easier to reason about.
- **"How does auth work?"** — stateless JWTs: the signed token carries the user
  ID, so the server verifies the signature instead of storing sessions.
- **"How are tasks kept private?"** — every query is scoped with `WHERE user_id = ?`;
  one user physically cannot read another's rows.
- **"How would you scale it?"** — swap SQLite for Postgres behind the same
  `TaskRepository` interface; the handlers don't change. (That's the payoff of
  the interface boundary.)
- **"How do you test handlers?"** — `httptest` fires requests at the router with no
  real network, against a temp database — fast and deterministic.

---

## 🎯 Final exercises — make it *yours*

1. **Ship it to GitHub.** Initialize a repo, commit, and push the whole `golang/`
   course (or just `taskflow/`). A green commit history is itself a signal.
2. **Add one feature** end-to-end and write its test: e.g. task **priorities**,
   **due dates**, or a `GET /tasks?done=true` filter. Touch every layer
   (model → store → handler → test).
3. **Build and run the Docker image** locally; confirm the API works in a container.
4. **Write your CV bullets** for TaskFlow using the template above, customized to
   the feature *you* added.
5. **Keep learning:** explore `log/slog` (structured logging), add a CI workflow
   (`.github/workflows`) that runs `go test ./...`, or swap in PostgreSQL.

---

## ✅ Session 20 Checklist

- [ ] My server shuts down gracefully on SIGINT/SIGTERM
- [ ] I have a multi-stage Dockerfile producing a small static image
- [ ] I understand why `CGO_ENABLED=0` works here (pure-Go SQLite driver)
- [ ] I have `.dockerignore` and `.gitignore`
- [ ] I have a Makefile with run/test/build targets
- [ ] My README has features, stack, architecture, run, API, and test sections
- [ ] I can describe TaskFlow in CV bullets with specifics and impact
- [ ] I can answer the "why" interview questions about my design
- [ ] I pushed the project to GitHub

**Previous:** [← Session 19](session-19.md)

---

## 🎓 You finished the course!

You went from `package main` to a containerized, authenticated, tested REST API.
You now know Go's syntax, its data structures, interfaces, errors, concurrency,
the standard library, testing, and HTTP — and you have a **real project** to prove
it. That's a junior-to-mid Go backend developer skill set.

**Where to go next:** build *another* small project from scratch (a URL shortener,
a expense tracker API) without looking at notes — that's how this becomes muscle
memory. Then read real Go code: the standard library, and projects like `chi`,
`sqlc`, or the source of tools you use. Welcome to the Go community. 🐹
