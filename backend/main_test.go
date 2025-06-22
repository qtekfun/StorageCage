// backend/main_test.go
package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestHomeHandler will test the handler for our root path.
func TestHomeHandler(t *testing.T) {
	// The expected response body from our handler.
	// We write it here so we have a single source of truth for the test.
	expectedBody := "StorageCage API server is running! 👋"

	// Create a new request to the root path "/".
	// We don't need to provide a body for a GET request, so we use nil.
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Create a ResponseRecorder, which is a special type that acts like a
	// ResponseWriter but records the results for testing.
	rr := httptest.NewRecorder()

	// Create an http.HandlerFunc from our homeHandler.
	// This allows us to test the handler directly without starting a full server.
	handler := http.HandlerFunc(homeHandler)

	// Serve the HTTP request to our ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check if the status code is what we expect (200 OK).
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check if the response body is what we expect.
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

// TestListFilesHandler will test our future /api/v1/files endpoint.
func TestListFilesHandler(t *testing.T) {
	// This is the exact JSON output we expect from our API.
	// Using raw string literals (`) makes it easy to write multi-line strings.
	expectedBody := `[{"id":"1","name":"file1.txt","size":1024},{"id":"2","name":"image.jpg","size":5242880}]`

	// We create a request to our API endpoint.
	req := httptest.NewRequest(http.MethodGet, "/api/v1/files", nil)

	// A ResponseRecorder to capture the response.
	rr := httptest.NewRecorder()

	// We need a router to route the request to the correct handler.
	// We'll create it here for the test.
	// This shows the power of having the newRouter() function!
	router := newRouter()
	router.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the Content-Type header.
	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			ctype, "application/json")
	}

	// Check the response body. We trim whitespace from the response body for a robust comparison.
    // NOTE: In a real-world scenario with complex JSON, it's better to unmarshal
    // the response into a struct and compare the struct fields. For now, this is fine.
	if strings.TrimSpace(rr.Body.String()) != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}