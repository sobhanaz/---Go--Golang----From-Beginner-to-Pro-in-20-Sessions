// Session 14 — strconv: convert between strings and numbers.
// Run:  go run examples/session14/strconv/strconv.go
package main

import (
	"fmt"
	"strconv"
)

func main() {
	// string -> int (Atoi = "ASCII to integer"). Returns a value AND an error.
	n, err := strconv.Atoi("42")
	if err != nil {
		fmt.Println("parse error:", err)
	} else {
		fmt.Println("Atoi:", n+1) // 43
	}

	// A failing parse: handle the error, don't ignore it.
	if _, err := strconv.Atoi("not-a-number"); err != nil {
		fmt.Println("expected error:", err)
	}

	// int -> string (Itoa = "integer to ASCII").
	s := strconv.Itoa(2026)
	fmt.Println("Itoa:", s+"!")

	// string -> float64.
	f, _ := strconv.ParseFloat("3.14", 64)
	fmt.Printf("ParseFloat: %.4f\n", f*2)

	// string -> bool.
	b, _ := strconv.ParseBool("true")
	fmt.Println("ParseBool:", b)

	// Format a float to a string with control over precision.
	fmt.Println("FormatFloat:", strconv.FormatFloat(3.14159, 'f', 2, 64)) // 3.14
}
