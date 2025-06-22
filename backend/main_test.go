// backend/main_test.go
package main

import (
	"bytes" // Needed to create an in-memory buffer for the request body
	"encoding/json"
	"io"             // Needed for io.Copy
	"mime/multipart" // The key package for building multipart/form-data requests
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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

func TestUploadFileHandler(t *testing.T) {
	// 1. Setup: Create a temp directory and the router config
	tempDir := t.TempDir()
	config := AppConfig{StorageDir: tempDir}
	router := newRouter(&config)

	// 2. Create the multipart request body
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	// The field name "file" is important, our handler will look for it.
	part, err := writer.CreateFormFile("file", "test-upload.txt")
	if err != nil {
		t.Fatal(err)
	}
	// Write the "file content" to the multipart request
	_, err = io.WriteString(part, "this is a test file content")
	if err != nil {
		t.Fatal(err)
	}
	writer.Close() // This finalizes the request body

	// 3. Create the HTTP request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/files", body)
	// Set the Content-Type header, which includes the multipart boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 4. Perform the request
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// 5. Assertions
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Check that the file was actually created on disk
	filePath := filepath.Join(tempDir, "test-upload.txt")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("handler did not create file on disk: %v", err)
	}
	if string(content) != "this is a test file content" {
		t.Errorf("file content is incorrect")
	}

	// Check the JSON response
	var fileInfo FileInfo
	if err := json.NewDecoder(rr.Body).Decode(&fileInfo); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	if fileInfo.Name != "test-upload.txt" {
		t.Errorf("response file name is incorrect: got %v want %v", fileInfo.Name, "test-upload.txt")
	}
}

func TestDeleteFileHandler(t *testing.T) {
	// 1. Setup: Create a temp directory and a dummy file to delete.
	tempDir := t.TempDir()
	dummyFileName := "file-to-delete.txt"
	dummyFilePath := filepath.Join(tempDir, dummyFileName)
	if err := os.WriteFile(dummyFilePath, []byte("delete me"), 0666); err != nil {
		t.Fatalf("Failed to create dummy file for deletion test: %v", err)
	}

	config := AppConfig{StorageDir: tempDir}
	router := newRouter(&config)

	// 2. Create the request to the DELETE endpoint.
	// The file name is passed as a URL parameter.
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/files/"+dummyFileName, nil)
	rr := httptest.NewRecorder()

	// 3. Perform the request.
	router.ServeHTTP(rr, req)

	// 4. Assertions
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// 4a. Check the success message in the response body.
	expected := `{"message":"file deleted successfully"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// 4b. The most important check: verify the file NO LONGER EXISTS on disk.
	if _, err := os.Stat(dummyFilePath); !os.IsNotExist(err) {
		t.Errorf("handler did not delete the file from disk, it still exists")
	}
}
