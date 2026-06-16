> 🌐 **Language / زبان:** English (this file) · [فارسی](session-19.fa.md)

# Session 19 — Auth, Middleware & Config 🔐

**Goal (1 hour):** Turn TaskFlow from a toy into a real backend. You'll add
**configuration** from the environment, **middleware** (logging, panic recovery,
auth), and **JWT authentication** with password hashing — and make every task
**scoped to its owner** so users only see their own data. This is the session
that makes your project genuinely portfolio-worthy.

> **Recap from Session 18:** TaskFlow has a database-backed, layered CRUD API
> (models → store → api). Now we secure it and make it production-shaped.

> 📂 All code lives in [`taskflow/`](../taskflow/). Run the tests with
> `cd taskflow && go test ./...`, or run the server with `go run .`.

---

## 1. Configuration from the environment (10 min)

Hardcoding ports, database paths, and secrets is a mistake — real apps read them
from the environment so the same binary runs in dev, staging, and production
unchanged. We built a tiny [`config`](../taskflow/internal/config/config.go) package:

```go
type Config struct {
    Addr        string // TASKFLOW_ADDR   (default ":8080")
    DatabaseDSN string // TASKFLOW_DB     (default "taskflow.db")
    JWTSecret   string // TASKFLOW_JWT_SECRET
}

func Load() Config { /* os.Getenv with fallbacks */ }
```

```bash
TASKFLOW_ADDR=:9000 TASKFLOW_JWT_SECRET=super-secret go run .
```

> 🔑 **Never commit real secrets.** The default JWT secret here is for local dev
> only. In production you inject `TASKFLOW_JWT_SECRET` via the environment (or a
> secrets manager). This is a point worth making in an interview.

---

## 2. Password hashing & JWTs (20 min)

Two new third-party packages (added with `go get`):

- `golang.org/x/crypto/bcrypt` — hash passwords.
- `github.com/golang-jwt/jwt/v5` — issue and verify tokens.

### Hashing passwords — never store plaintext

[`auth/password.go`](../taskflow/internal/auth/password.go):

```go
func HashPassword(password string) (string, error) {
    b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(b), err
}

func CheckPassword(hash, password string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
```

> 🔑 **bcrypt is deliberately slow and salts automatically.** You store only the
> hash; you can never get the password back. At login you *compare*, you never
> decrypt. Storing plaintext passwords is the #1 security sin — bcrypt is the fix.

### JSON Web Tokens — stateless auth

A **JWT** is a signed token the client sends on every request to prove who it is.
[`auth/jwt.go`](../taskflow/internal/auth/jwt.go) puts the user's ID in the token's
`Subject` and signs it with the secret:

```go
func GenerateToken(secret string, userID int64) (string, error) {
    claims := jwt.RegisteredClaims{
        Subject:   strconv.FormatInt(userID, 10),
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
```

`ParseToken` reverses this: it verifies the signature (rejecting unexpected
signing algorithms — a real security check), confirms the token isn't expired,
and returns the user ID. Because the token is **signed**, the server doesn't need
to store sessions — it can trust a valid token. That's "stateless" auth.

---

## 3. Middleware — cross-cutting behavior (15 min)

A **middleware** wraps a handler to add behavior, with the signature:

```go
func(next http.Handler) http.Handler
```

It returns a new handler that does something *before/after* calling `next`. We
built three in [`api/middleware.go`](../taskflow/internal/api/middleware.go):

- **`Logging`** — records `method path -> status (duration)` for every request.
  It wraps the `ResponseWriter` in a `statusRecorder` to capture the status code.
- **`Recovery`** — `defer`s a `recover()` (Session 11!) so a panic in any handler
  returns a clean 500 instead of crashing the whole server.
- **`Auth`** — reads the `Authorization: Bearer <token>` header, verifies the JWT,
  and stores the user ID in the request **context**. No token → `401`.

### Passing data through `context`

`Auth` puts the authenticated user's ID into the request context so handlers can
read it:

```go
ctx := context.WithValue(r.Context(), userIDKey, userID)
next.ServeHTTP(w, r.WithContext(ctx))
```

```go
func userIDFromContext(r *http.Request) int64 {
    id, _ := r.Context().Value(userIDKey).(int64)
    return id
}
```

> 🔑 **Use an unexported key type** (`type ctxKey string`) for context keys so
> other packages can't accidentally collide with yours. This is the idiomatic
> pattern for request-scoped values.

### Applying middleware

Global middleware wraps the whole mux; per-route middleware wraps individual
handlers. In [`api/server.go`](../taskflow/internal/api/server.go):

```go
// Public routes — no auth.
mux.HandleFunc("POST /auth/register", s.handleRegister)
mux.HandleFunc("POST /auth/login", s.handleLogin)

// Protected routes — wrapped with Auth.
mux.Handle("GET /tasks", s.Auth(http.HandlerFunc(s.handleListTasks)))
// ...

// Global chain: Recovery is outermost so it catches everything.
return Recovery(Logging(mux))
```

---

## 4. User-scoped data — the realistic part (15 min)

A real multi-user app must isolate data. We made **every task belong to a user**:

- The `tasks` table got a `user_id` column with a foreign key to `users`.
- Every [`TaskStore`](../taskflow/internal/store/task_store.go) method now takes a
  `userID` and filters by it: `WHERE id = ? AND user_id = ?`.
- Handlers read the ID from the context (set by `Auth`) and pass it down.

The payoff is real isolation, proven by a test: when **Bob** asks for **Alice's**
task by ID, he gets a `404` — not because it doesn't exist, but because it isn't
*his*. The SQL `AND user_id = ?` makes leaking another user's data impossible.

### The auth handlers tie it together

[`api/auth.go`](../taskflow/internal/api/auth.go):
- **`POST /auth/register`** — validate input, hash the password, create the user,
  return a JWT. Duplicate email → `409 Conflict`.
- **`POST /auth/login`** — look up the user, `CheckPassword`, return a JWT. Wrong
  email *or* password → the **same** generic `401` message (so you don't reveal
  which emails are registered — a subtle security best practice).

Notice the login response embeds the `User`, but the password hash never appears
because of its `json:"-"` tag (Session 15). Secrets stay server-side.

---

## 5. See it all work (5 min)

```bash
cd taskflow
go test ./...                       # all auth + scoping tests pass

go run .                            # start the server, then in another terminal:
curl localhost:8080/health
curl -X POST localhost:8080/auth/register -d '{"email":"me@example.com","password":"secret123"}'
# copy the "token" from the response, then:
TOKEN=...
curl -X POST localhost:8080/tasks -H "Authorization: Bearer $TOKEN" -d '{"title":"Ship it"}'
curl localhost:8080/tasks -H "Authorization: Bearer $TOKEN"
```

The test suite ([`api/tasks_test.go`](../taskflow/internal/api/tasks_test.go)) covers
the whole story: public health, 401 without a token, full CRUD with auth,
**per-user isolation**, the login flow, and duplicate-email rejection — all with
`httptest` and a real (temporary) database.

---

## 🎯 Exercises (do these before Session 20!)

Work in the `taskflow/` project:

1. **`GET /me`:** Add a protected endpoint that returns the current user's info
   (id, email). You'll need a `UserStore.GetByID` method. Don't leak the hash.
2. **Stronger validation:** Reject registration when the email has no `@`, and
   require passwords ≥ 8 chars. Return clear `400` messages.
3. **CORS middleware:** Add a `CORS` middleware that sets
   `Access-Control-Allow-Origin: *` and handles `OPTIONS` preflight requests.
4. **Token expiry test:** Write a test that generates a token with a past expiry
   (tweak `GenerateToken` or craft claims) and asserts the API returns `401`.
5. **Auth unit tests:** Add a `password_test.go` and `jwt_test.go` in the `auth`
   package: hash→check round-trip, and generate→parse round-trip plus a
   tampered-token failure.

---

## ✅ Session 19 Checklist

- [ ] I load configuration from environment variables with sensible defaults
- [ ] I hash passwords with bcrypt and never store plaintext
- [ ] I can generate and verify a JWT, and reject expired/tampered tokens
- [ ] I understand the `func(next http.Handler) http.Handler` middleware shape
- [ ] I built logging, recovery, and auth middleware
- [ ] I pass the user ID through `context` with an unexported key type
- [ ] My task data is scoped per user (`WHERE user_id = ?`)
- [ ] Login uses a generic error for both wrong-email and wrong-password
- [ ] I completed all 5 exercises

**Previous:** [← Session 18](session-18.md) · **Next:** [Session 20 — Polish & Ship →](session-20.md)
