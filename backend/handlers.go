package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// FileInfo represents metadata for a file in our system.
type FileInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

// homeHandler handles requests to the root ("/") path.
func (cfg *AppConfig) homeHandler(w http.ResponseWriter, r *http.Request) { // <-- Cambio
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "StorageCage API server is running! 👋")
}

func (cfg *AppConfig) listFilesHandler(w http.ResponseWriter, r *http.Request) {
	// Use the os.ReadDir function to read the directory specified in our config.
	dirEntries, err := os.ReadDir(cfg.StorageDir)
	if err != nil {
		// If the directory doesn't exist or we can't read it, return an internal server error.
		log.Printf("ERROR: could not read storage directory: %v", err)
		http.Error(w, "Could not list files", http.StatusInternalServerError)
		return
	}

	// Create a slice to hold our file information.
	var files []FileInfo

	// Iterate over each entry in the directory.
	for _, entry := range dirEntries {
		// Skip subdirectories.
		if entry.IsDir() {
			continue
		}

		// Get detailed file info (which includes size).
		info, err := entry.Info()
		if err != nil {
			// If we can't get info for a specific file, log it and skip it.
			log.Printf("ERROR: could not get file info for %s: %v", entry.Name(), err)
			continue
		}

		// Create our FileInfo struct and append it to the slice.
		fileData := FileInfo{
			// For now, ID is the same as the name. We can improve this later.
			ID:   info.Name(),
			Name: info.Name(),
			Size: info.Size(),
		}
		files = append(files, fileData)
	}

	// The rest is the same: set header and encode the slice to JSON.
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(files); err != nil {
		log.Printf("ERROR: Error encoding files to JSON: %v", err)
	}
}
