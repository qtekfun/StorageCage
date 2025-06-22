// backend/main_test.go
package main

import (
	"net/http"
	"net/http/httptest"
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