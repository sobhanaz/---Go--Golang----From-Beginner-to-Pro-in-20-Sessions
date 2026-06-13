// Session 14 — the time package: now, durations, formatting, parsing.
// Run:  go run examples/session14/timepkg/timepkg.go
package main

import (
	"fmt"
	"time"
)

func main() {
	// A fixed time so this example's output is stable (normally: time.Now()).
	t := time.Date(2026, time.June, 13, 15, 4, 5, 0, time.UTC)
	fmt.Println("time:", t)

	// Extract parts.
	fmt.Printf("year=%d month=%s day=%d hour=%d\n",
		t.Year(), t.Month(), t.Day(), t.Hour())

	// Durations: add/subtract time.
	later := t.Add(48 * time.Hour)
	fmt.Println("48h later:", later.Format("2006-01-02"))

	// Difference between two times is a Duration.
	diff := later.Sub(t)
	fmt.Println("difference in hours:", diff.Hours())

	// Formatting uses Go's REFERENCE date: Mon Jan 2 15:04:05 MST 2006.
	// You write the layout AS THAT date; Go fills in your value.
	fmt.Println("formatted:", t.Format("2006-01-02 15:04:05"))
	fmt.Println("friendly: ", t.Format("Mon, 02 Jan 2006"))

	// Parsing a string into a time (same reference-layout idea).
	parsed, err := time.Parse("2006-01-02", "2025-12-25")
	if err == nil {
		fmt.Println("parsed:", parsed.Format("Monday, Jan 2"))
	}
}
