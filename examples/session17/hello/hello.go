// Session 17 — the smallest HTTP server.
// Run:  go run examples/session17/hello/hello.go
// Then visit http://localhost:8080/ and http://localhost:8080/hello in a browser,
// or in another terminal:  curl localhost:8080/hello
//
// Press Ctrl+C to stop the server.
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// A ServeMux is a request router: it maps URL patterns to handlers.
	mux := http.NewServeMux()

	// Each handler receives a ResponseWriter (to write the reply) and a
	// *Request (the incoming request). This is the core signature in Go web dev.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to your first Go server!")
	})

	// Go 1.22+ lets you specify the METHOD and PATH together.
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		// Query parameters: /hello?name=Sobhan
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "stranger"
		}
		fmt.Fprintf(w, "Hello, %s!\n", name)
	})

	addr := ":8080"
	log.Printf("listening on http://localhost%s", addr)
	// ListenAndServe BLOCKS, serving requests until the program is stopped.
	// It only returns if there's an error (e.g. the port is already in use).
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
