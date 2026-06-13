> 🌐 **Language / زبان:** English (this file) · [فارسی](README.fa.md)

# 🐹 Go (Golang) — From Beginner to Pro in 20 Sessions

Welcome! This is your complete, self-paced course to learn the Go programming
language from zero to a job-ready level. Each session is designed for **~1 hour**.
By the end you'll have built a **real, portfolio-worthy project** you can put on
your CV and discuss in interviews.

> **Your setup is verified:** Go 1.25.5, git, and VS Code are installed.
> Module name for this course: `golearn`.

---

## How to use this course

1. Do **one session per day** (or at your own pace). Don't rush — type every example yourself.
2. Each session lives in [`sessions/`](sessions/) as its own README. Read it top to bottom.
3. Runnable code for each session lives in [`examples/`](examples/). Run it, break it, fix it.
4. Every session ends with **exercises**. Do them before moving on — that's where learning sticks.
5. Keep this repo on GitHub. By session 20 it *is* your portfolio.

### The two commands you'll use constantly

```bash
go run path/to/file.go      # compile + run in one step (for learning)
go build ./...              # compile everything into a binary
go test ./...               # run all tests
go fmt ./...                # auto-format your code (always do this)
```

---

## 📚 The 20-Session Roadmap

### Part 1 — Foundations (Sessions 1–5)
| # | Session | What you'll learn |
|---|---------|-------------------|
| 01 | [Hello, Go](sessions/session-01.md) | Install check, your first program, `go run`/`go build`, packages, `gofmt` |
| 02 | Variables & Types | `var`, `:=`, constants, basic types, zero values, type conversion |
| 03 | Control Flow | `if`, `switch`, `iota`, operators, boolean logic |
| 04 | Loops & Functions | the `for` loop (Go's only loop), defining and calling functions |
| 05 | Functions Deep Dive | multiple returns, named returns, variadic, closures, `defer` |

### Part 2 — Data Structures (Sessions 6–9)
| # | Session | What you'll learn |
|---|---------|-------------------|
| 06 | Arrays & Slices | fixed arrays vs. dynamic slices, `append`, `make`, slicing, `copy` |
| 07 | Maps, Strings & Runes | key/value maps, string internals, UTF-8, `rune`, `byte` |
| 08 | Structs & Methods | custom types, methods, value vs. pointer receivers, embedding |
| 09 | Pointers | what pointers are, `&` and `*`, when and why to use them |

### Part 3 — The Go Way (Sessions 10–13)
| # | Session | What you'll learn |
|---|---------|-------------------|
| 10 | Interfaces | Go's superpower: implicit interfaces, polymorphism, `any` |
| 11 | Errors | idiomatic error handling, custom errors, `errors.Is/As`, `panic`/`recover` |
| 12 | Concurrency I | goroutines, channels, the `go` keyword |
| 13 | Concurrency II | `select`, `sync.WaitGroup`, `Mutex`, `context`, worker pools |

### Part 4 — Real-World Go (Sessions 14–17)
| # | Session | What you'll learn |
|---|---------|-------------------|
| 14 | Standard Library Tour | `fmt`, `strings`, `strconv`, `time`, `os`, `sort` |
| 15 | Files, JSON & Encoding | reading/writing files, `encoding/json`, struct tags |
| 16 | Testing | unit tests, table-driven tests, benchmarks, coverage |
| 17 | HTTP Servers | `net/http`, handlers, routing, JSON APIs |

### Part 5 — The Portfolio Project (Sessions 18–20)
| # | Session | What you'll learn |
|---|---------|-------------------|
| 18 | REST API + Database | project layout, connecting to a database, CRUD endpoints |
| 19 | Auth, Middleware & Config | JWT auth, middleware, env config, structured logging |
| 20 | Polish & Ship | Docker, tests, README, deploy — and your **CV bullet points** |

---

## 🏁 Final Project: "TaskFlow" — a Task Manager REST API

By session 20 you will have built **TaskFlow**, a production-style backend:

- ✅ RESTful JSON API built with Go's standard library + a lightweight router
- ✅ PostgreSQL (or SQLite) persistence with proper migrations
- ✅ User registration & login with **JWT authentication**
- ✅ Middleware (logging, auth, recovery, CORS)
- ✅ Clean, layered architecture (handler → service → repository)
- ✅ Unit and integration tests
- ✅ Dockerized and ready to deploy
- ✅ A polished README so recruiters can run it in one command

This is the kind of project that gets you a **Go backend developer** interview.

---

## ✅ Progress Tracker

- [x] Session 01 — Hello, Go
- [x] Session 02 — Variables & Types
- [x] Session 03 — Control Flow
- [x] Session 04 — Loops & Functions
- [x] Session 05 — Functions Deep Dive
- [x] Session 06 — Arrays & Slices
- [x] Session 07 — Maps, Strings & Runes
- [x] Session 08 — Structs & Methods
- [x] Session 09 — Pointers
- [x] Session 10 — Interfaces
- [x] Session 11 — Errors
- [x] Session 12 — Concurrency I
- [x] Session 13 — Concurrency II
- [x] Session 14 — Standard Library Tour
- [ ] Session 15 — Files, JSON & Encoding
- [ ] Session 16 — Testing
- [ ] Session 17 — HTTP Servers
- [ ] Session 18 — REST API + Database
- [ ] Session 19 — Auth, Middleware & Config
- [ ] Session 20 — Polish & Ship

---

*Course built for you with Go 1.25.5. Start with [Session 01](sessions/session-01.md). Happy coding! 🚀*
