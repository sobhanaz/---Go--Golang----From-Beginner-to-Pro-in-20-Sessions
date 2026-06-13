// Session 14 — the os package: args, environment, and exit codes.
// Run:  go run examples/session14/ospkg/ospkg.go hello world
package main

import (
	"fmt"
	"os"
)

func main() {
	// os.Args holds the command-line arguments.
	// os.Args[0] is the program path; the rest are user arguments.
	fmt.Println("program:", os.Args[0])
	fmt.Println("arg count (excluding program):", len(os.Args)-1)
	for i, arg := range os.Args[1:] {
		fmt.Printf("  arg %d: %s\n", i, arg)
	}

	// Environment variables.
	home := os.Getenv("HOME")
	fmt.Println("HOME:", home)

	// LookupEnv distinguishes "empty" from "not set" (comma-ok style).
	if v, ok := os.LookupEnv("DEFINITELY_NOT_SET_12345"); !ok {
		fmt.Println("env var not set; would default to something")
	} else {
		fmt.Println("got:", v)
	}

	// os.Exit(code) ends the program immediately with a status code.
	// 0 = success, non-zero = error. (Deferred funcs do NOT run on os.Exit!)
	// We don't call it here so the program ends normally.
	fmt.Println("done")
}
