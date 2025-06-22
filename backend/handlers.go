package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

// FileInfo represents metadata for a file in our system.
type FileInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type StatusResponse struct {
	Message string `json:"message"`
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

func (cfg *AppConfig) uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Parse the multipart form. 10 << 20 specifies a maximum
	// upload of 10 MB. It's a good idea to have a limit.
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	// 2. Retrieve the file from the form data.
	// "file" is the key that the client is expected to use for the file upload.
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file key", http.StatusBadRequest)
		return
	}
	defer file.Close() // Crucial to close the file when we're done

	// 3. Create the destination file on the server.
	dstPath := filepath.Join(cfg.StorageDir, handler.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Could not create file on server", http.StatusInternalServerError)
		return
	}
	defer dst.Close() // Also crucial

	// 4. Copy the uploaded file's content to the destination file.
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Could not save file content", http.StatusInternalServerError)
		return
	}

	// 5. Respond with the details of the created file.
	fileInfo := FileInfo{
		ID:   handler.Filename,
		Name: handler.Filename,
		Size: handler.Size,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Set the status code to 201
	json.NewEncoder(w).Encode(fileInfo)
}

func (cfg *AppConfig) deleteFileHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get the file name from the URL parameter.
	// Chi makes this easy with chi.URLParam(request, "paramName").
	fileName := chi.URLParam(r, "fileName")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	// 2. IMPORTANT - Security: Sanitize the file name to prevent path traversal attacks.
	// For example, a request to /api/v1/files/../../importantsystemfile should not be allowed.
	// filepath.Base() strips all directory information, ensuring we only get a file name.
	fileName = filepath.Base(fileName)

	// 3. Construct the full path to the file.
	filePath := filepath.Join(cfg.StorageDir, fileName)

	// 4. Attempt to remove the file from disk.
	err := os.Remove(filePath)
	if err != nil {
		// If the error is that the file doesn't exist, return a 404 Not Found.
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		// For any other error (e.g., permissions), return a 500 Internal Server Error.
		log.Printf("ERROR: could not delete file %s: %v", filePath, err)
		http.Error(w, "Could not delete file", http.StatusInternalServerError)
		return
	}

	// 5. Respond with a success message.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(StatusResponse{Message: "file deleted successfully"})
}
