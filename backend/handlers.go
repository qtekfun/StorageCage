package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// FileInfo represents metadata for a file in our system.
type FileInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

// homeHandler handles requests to the root ("/") path.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "StorageCage API server is running! 👋")
}

// listFilesHandler creates a mock list of files and returns it as JSON.
func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	files := []FileInfo{
		{ID: "1", Name: "file1.txt", Size: 1024},
		{ID: "2", Name: "image.jpg", Size: 5242880},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(files); err != nil {
		log.Printf("Error encoding files to JSON: %v", err)
	}
}
