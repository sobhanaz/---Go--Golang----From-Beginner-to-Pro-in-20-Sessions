// Session 17 — testing HTTP handlers WITHOUT starting a real server.
// net/http/httptest lets you send fake requests and inspect the response.
// Run:  go test -v ./examples/session17/jsonapi/
package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealth(t *testing.T) {
	srv := NewServer()
	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder() // records what the handler writes

	srv.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"status":"ok"`) {
		t.Errorf("unexpected body: %s", rec.Body.String())
	}
}

func TestCreateAndGetTask(t *testing.T) {
	srv := NewServer()

	// 1. Create a task via POST.
	body := strings.NewReader(`{"title":"Learn Go"}`)
	req := httptest.NewRequest("POST", "/tasks", body)
	rec := httptest.NewRecorder()
	srv.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d; want 201", rec.Code)
	}
	var created Task
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	if created.ID != 1 || created.Title != "Learn Go" {
		t.Errorf("unexpected task: %+v", created)
	}

	// 2. Fetch it back via GET /tasks/1.
	req2 := httptest.NewRequest("GET", "/tasks/1", nil)
	rec2 := httptest.NewRecorder()
	srv.routes().ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("get status = %d; want 200", rec2.Code)
	}
}

func TestCreateTaskValidation(t *testing.T) {
	srv := NewServer()
	// Empty title should be rejected with 400.
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":""}`))
	rec := httptest.NewRecorder()
	srv.routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d; want 400", rec.Code)
	}
}

func TestGetMissingTask(t *testing.T) {
	srv := NewServer()
	req := httptest.NewRequest("GET", "/tasks/999", nil)
	rec := httptest.NewRecorder()
	srv.routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d; want 404", rec.Code)
	}
}
