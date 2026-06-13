> 🌐 **Language / زبان:** English (this file) · [فارسی](session-14.fa.md)

# Session 14 — Standard Library Tour 🧰

**Goal (1 hour):** Go ships with a famously rich **standard library** — batteries
included. You'll tour the packages you'll reach for daily: `strconv` (string↔number),
`time` (dates & durations), `sort` (ordering), and `os` (program environment).
Knowing the standard library well is a big part of being productive in Go.

> **Recap from Session 13:** you've finished Go's core concepts. From here it's
> about *building* — and the standard library is your toolbox.

> 💡 **How to explore any package:** run `go doc strconv` or `go doc strconv.Atoi`
> in your terminal, or browse [pkg.go.dev](https://pkg.go.dev). Get comfortable
> reading docs — it's a daily skill for a Go developer.

---

## 1. `strconv` — strings ↔ numbers (15 min)

Recall from Session 02 that Go won't convert a string to a number for you, and
`string(65)` gives `"A"`, not `"65"`. The `strconv` package does real conversions:

```go
n, err := strconv.Atoi("42")        // string -> int  ("ASCII to integer")
s := strconv.Itoa(2026)             // int -> string  ("integer to ASCII")
f, err := strconv.ParseFloat("3.14", 64) // string -> float64
b, err := strconv.ParseBool("true")      // string -> bool
```

> 🔑 **They return an error**, because the input might not be a valid number.
> Always check it — this is the `value, err` pattern from Session 05 in action.
> Parsing user input (form fields, query params, config) is where you'll use this
> constantly in the final project.

Run [`examples/session14/strconv/strconv.go`](../examples/session14/strconv/strconv.go).

---

## 2. `time` — dates, durations & formatting (20 min)

```go
now := time.Now()                  // current time
t := time.Date(2026, time.June, 13, 15, 4, 5, 0, time.UTC) // a specific time
```

### Durations

A `time.Duration` represents elapsed time. Build them by multiplying a unit:

```go
2 * time.Hour
500 * time.Millisecond
later := now.Add(48 * time.Hour)   // add time
diff  := later.Sub(now)            // subtract two times -> a Duration
fmt.Println(diff.Hours())          // 48
```

You already used durations in concurrency (`time.After`, `time.Sleep`).

### Formatting & parsing — Go's unique reference date

This trips everyone up, so learn it once. Instead of codes like `%Y-%m-%d`, Go
uses a **reference date** as the layout:

```
Mon Jan 2 15:04:05 MST 2006
```

(Read it as 1-2-3-4-5-6: 15h=3pm, Jan=1, 2nd, …, 2006.) You write the layout *as
that exact date*, and Go substitutes your value:

```go
t.Format("2006-01-02")          // -> "2026-06-13"
t.Format("Mon, 02 Jan 2006")    // -> "Sat, 13 Jun 2026"

parsed, err := time.Parse("2006-01-02", "2025-12-25") // string -> time
```

> 🔑 **Memorize `2006-01-02 15:04:05`** — it's the most common layout and the key
> to all date formatting in Go.

Run [`examples/session14/timepkg/timepkg.go`](../examples/session14/timepkg/timepkg.go).

---

## 3. `sort` — ordering data (15 min)

```go
sort.Ints(nums)       // sort an []int ascending, in place
sort.Strings(words)   // sort an []string
sort.Float64s(fs)     // sort an []float64
```

### Sorting structs (and anything custom) with `sort.Slice`

Pass a "less" function describing the order you want:

```go
sort.Slice(people, func(i, j int) bool {
    return people[i].Age < people[j].Age   // ascending by age
})
```

Return `people[i].Age > people[j].Age` for descending. The function answers
"should element i come before element j?"

### The idiom for ordered map output

Maps iterate randomly (Session 07). To print a map in a stable order, collect the
keys into a slice, sort them, then iterate:

```go
keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
sort.Strings(keys)
for _, k := range keys {
    fmt.Println(k, m[k])
}
```

Run [`examples/session14/sortpkg/sortpkg.go`](../examples/session14/sortpkg/sortpkg.go).

---

## 4. `os` — talking to the operating system (10 min)

```go
os.Args            // []string of command-line arguments (Args[0] = program path)
os.Getenv("HOME")  // read an environment variable ("" if not set)
os.LookupEnv("X")  // value, ok — distinguishes "empty" from "not set"
os.Exit(1)         // end the program NOW with a status code (0 = success)
```

> ⚠️ **`os.Exit` skips deferred functions!** If you've `defer`red cleanup, it
> won't run on `os.Exit`. Prefer returning errors up to `main` and exiting there.

Environment variables are how real apps get configuration (database URLs, ports,
secrets) — you'll read them in the final project's config layer.

Run it *with arguments*:

```bash
go run examples/session14/ospkg/ospkg.go hello world
```

> 📦 **Other standard-library gems** worth knowing exist (we'll meet several
> later): `math`, `math/rand`, `bufio`, `regexp`, `encoding/json` (Session 15),
> `net/http` (Session 17), `testing` (Session 16). The standard library is huge
> and high-quality — reach for it before adding a third-party dependency.

---

## 🎯 Exercises (do these before Session 15!)

Create `examples/session14/practice/practice.go`:

1. **Safe parse:** Write `func parseIntOrDefault(s string, def int) int` that uses
   `strconv.Atoi` and returns `def` if parsing fails. Test with `"100"` and `"oops"`.
2. **Age from year:** Read a birth year as a string, convert it, and print the
   age given the current year is 2026. Handle a bad input gracefully.
3. **Countdown formatting:** Given `time.Date(2027, 1, 1, ...)`, compute the
   duration from a fixed "now" and print the number of days until then.
4. **Sort people:** Make a `[]struct{Name string; Score int}` and sort it by score
   **descending**, then print the leaderboard.
5. **Ordered map:** Build a `map[string]int` of word counts and print the entries
   sorted alphabetically by word.

---

## ✅ Session 14 Checklist

- [ ] I can convert strings↔numbers with `strconv` and handle the error
- [ ] I can get the current time, add durations, and diff two times
- [ ] I can format and parse dates using the `2006-01-02 15:04:05` reference layout
- [ ] I can sort slices of ints/strings and structs (`sort.Slice`)
- [ ] I can print a map in sorted key order
- [ ] I can read `os.Args` and environment variables
- [ ] I know `os.Exit` skips deferred functions
- [ ] I completed all 5 exercises

**Previous:** [← Session 13](session-13.md) · **Next:** [Session 15 — Files, JSON & Encoding →](session-15.md)
