// backend/main_test.go
package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestListFilesHandler_RealFiles(t *testing.T) {
	// 1. Create a temporary directory for our test files.
	// t.TempDir() automatically creates it and cleans it up when the test finishes.
	tempDir := t.TempDir()

	// 2. Create a dummy file inside the temporary directory.
	dummyFileName := "testfile.txt"
	dummyFileContent := "hello world"
	err := os.WriteFile(filepath.Join(tempDir, dummyFileName), []byte(dummyFileContent), 0666)
	if err != nil {
		t.Fatalf("Failed to create dummy file: %v", err)
	}

	// 3. Setup our config and router for the test.
	config := AppConfig{StorageDir: tempDir}
	router := newRouter(&config)

	// 4. Perform the request.
	req := httptest.NewRequest(http.MethodGet, "/api/v1/files", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// 5. Assertions
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Decode the JSON response
	var files []FileInfo
	if err := json.NewDecoder(rr.Body).Decode(&files); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Check if we got the correct number of files.
	if len(files) != 1 {
		t.Errorf("handler returned wrong number of files: got %v want %v", len(files), 1)
	}

	// Check the details of the file.
	if files[0].Name != dummyFileName {
		t.Errorf("handler returned wrong file name: got %v want %v", files[0].Name, dummyFileName)
	}
	if files[0].Size != int64(len(dummyFileContent)) {
		t.Errorf("handler returned wrong file size: got %v want %v", files[0].Size, len(dummyFileContent))
	}
}
