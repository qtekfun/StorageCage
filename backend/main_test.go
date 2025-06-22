package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	router := newRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "StorageCage API server is running! 👋"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestListFilesHandler(t *testing.T) {
	expectedBody := `[{"id":"1","name":"file1.txt","size":1024},{"id":"2","name":"image.jpg","size":5242880}]`
	req := httptest.NewRequest(http.MethodGet, "/api/v1/files", nil)
	rr := httptest.NewRecorder()
	router := newRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			ctype, "application/json")
	}

	if strings.TrimSpace(rr.Body.String()) != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}
