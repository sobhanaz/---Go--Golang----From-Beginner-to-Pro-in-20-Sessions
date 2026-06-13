> 🌐 **Language / زبان:** English (this file) · [فارسی](session-17.fa.md)

# Session 17 — HTTP Servers 🌐

**Goal (1 hour):** Build web servers with Go's standard library — no framework
required. You'll write handlers, route requests, read query/path parameters and
JSON bodies, and return JSON responses with proper status codes. This is the
exact foundation the final **TaskFlow** project is built on.

> **Recap from Session 16:** you can test code. In this session you'll *also* see
> how to test HTTP handlers with `httptest` — without ever starting a real server.

---

## 1. The smallest server (15 min)

Go's `net/http` package is a production-grade web server in the standard library
(it powers huge systems). The smallest server:

```go
func main() {
    mux := http.NewServeMux()   // a router: maps URL patterns -> handlers

    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello from Go!")
    })

    http.ListenAndServe(":8080", mux)   // BLOCKS, serving until stopped
}
```

The heart of Go web development is the **handler signature**:

```go
func(w http.ResponseWriter, r *http.Request)
```

- `w http.ResponseWriter` — you **write the response** into this (body, headers, status).
- `r *http.Request` — the **incoming request** (method, URL, headers, body).

Run it, then visit the URL or curl it:

```bash
go run examples/session17/hello/hello.go
# in another terminal:
curl localhost:8080/hello?name=Sobhan
```

(Press Ctrl+C to stop the server.) Run [`examples/session17/hello/hello.go`](../examples/session17/hello/hello.go).

---

## 2. Routing: methods, paths & parameters (15 min)

Since **Go 1.22**, the built-in router understands HTTP methods and path
wildcards — so you often don't need a third-party router at all:

```go
mux.HandleFunc("GET /tasks", listTasks)        // method + path
mux.HandleFunc("POST /tasks", createTask)
mux.HandleFunc("GET /tasks/{id}", getTask)     // {id} is a path wildcard
```

Reading the different kinds of input:

```go
// Path parameter:  GET /tasks/42  ->  r.PathValue("id") == "42"
id := r.PathValue("id")

// Query parameter:  GET /search?q=go  ->  r.URL.Query().Get("q")
q := r.URL.Query().Get("q")

// JSON request body (e.g. from a POST):
var input struct{ Title string `json:"title"` }
json.NewDecoder(r.Body).Decode(&input)
```

> 🔑 `r.PathValue("id")` returns a **string** — convert it with `strconv.Atoi`
> (Session 14) and handle the error if it's not a valid number.

---

## 3. Returning JSON properly (15 min)

A REST API responds with JSON and a meaningful **HTTP status code**. The clean
pattern is a small helper:

```go
func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")  // 1. set the content type
    w.WriteHeader(status)                                // 2. set the status code
    json.NewEncoder(w).Encode(v)                         // 3. write the JSON body
}
```

> ⚠️ **Order matters!** Set headers *first*, then call `WriteHeader(status)`, then
> write the body. Once you write the body, the status is locked in. Calling
> `WriteHeader` after writing always defaults to 200.

Common status codes you'll use:

| Code | Constant | Meaning |
|------|----------|---------|
| 200 | `http.StatusOK` | success (GET) |
| 201 | `http.StatusCreated` | created (POST) |
| 400 | `http.StatusBadRequest` | bad input from the client |
| 404 | `http.StatusNotFound` | resource doesn't exist |
| 500 | `http.StatusInternalServerError` | server-side failure |

Use the **named constants**, not raw numbers — they're clearer and self-documenting.

---

## 4. Structure for testability — the `Server` struct (15 min)

A professional pattern (and the one we'll use in TaskFlow): put your dependencies
(database, config, in-memory store) in a `Server` struct, make handlers
**methods** on it, and expose a `routes()` method that returns the router.

```go
type Server struct {
    mu    sync.Mutex
    tasks map[int]Task   // later: a real database
}

func (s *Server) handleListTasks(w http.ResponseWriter, r *http.Request) { ... }

func (s *Server) routes() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /tasks", s.handleListTasks)
    mux.HandleFunc("POST /tasks", s.handleCreateTask)
    return mux
}
```

Why this matters: it makes handlers **trivially testable**. The `net/http/httptest`
package lets you fire a fake request at `routes()` and inspect the response —
**no real server, no real port:**

```go
srv := NewServer()
req := httptest.NewRequest("GET", "/health", nil)
rec := httptest.NewRecorder()        // captures the response
srv.routes().ServeHTTP(rec, req)

if rec.Code != http.StatusOK { t.Fatalf("got %d", rec.Code) }
```

This combines Session 16 (testing) with this session, and it's how real Go APIs
are tested in CI. Study the full example and its tests:

```bash
go run examples/session17/jsonapi/jsonapi.go    # run the real server
go test -v ./examples/session17/jsonapi/        # test the handlers (no server)
```

See [`examples/session17/jsonapi/jsonapi.go`](../examples/session17/jsonapi/jsonapi.go) and
[`jsonapi_test.go`](../examples/session17/jsonapi/jsonapi_test.go).

> 📦 **Middleware preview:** a middleware is a function that wraps a handler to add
> cross-cutting behavior (logging, auth, recovery). We'll build these in Session 19.
> Their shape: `func(next http.Handler) http.Handler`.

---

## 🎯 Exercises (do these before Session 18!)

Work in `examples/session17/practice/`:

1. **Echo server:** A `GET /echo?msg=hi` endpoint that returns JSON
   `{"echo":"hi"}`. Default to `"nothing"` when `msg` is absent.
2. **Greeting with path param:** `GET /greet/{name}` returning
   `{"greeting":"Hello, <name>"}`.
3. **In-memory notes API:** `POST /notes` (create from JSON body) and
   `GET /notes` (list all). Use a `Server` struct with a slice or map.
4. **Status codes:** Make `POST /notes` return 400 if the body is missing a
   `text` field, and 201 on success.
5. **Test it:** Write `httptest`-based tests for your notes API covering create,
   list, and the 400 validation case. Run `go test`.

---

## ✅ Session 17 Checklist

- [ ] I can start a server with `http.NewServeMux` + `http.ListenAndServe`
- [ ] I understand the `func(w, r)` handler signature
- [ ] I can route by method and path, including `{id}` path wildcards
- [ ] I can read path params, query params, and a JSON body
- [ ] I can return JSON with the correct `Content-Type` and status code
- [ ] I know to set headers/status before writing the body
- [ ] I structure handlers as methods on a `Server` for testability
- [ ] I can test handlers with `httptest` (no real server)
- [ ] I completed all 5 exercises

**Previous:** [← Session 16](session-16.md) · **Next:** [Session 18 — REST API + Database →](session-18.md)

---

🎉 **Milestone:** Part 4 (Real-World Go) complete! You can now build tested JSON
web services. Next we begin **Part 5** — building the full TaskFlow portfolio
project, starting with database-backed storage.
