// backend/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
)

// FileInfo represents metadata for a file in our system.
// The `json:"..."` tags are used to control how the struct fields are
// named when encoded into JSON.
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
	// Create a slice of FileInfo structs with mock data.
	// This is the data we want to send as a response.
	files := []FileInfo{
		{ID: "1", Name: "file1.txt", Size: 1024},
		{ID: "2", Name: "image.jpg", Size: 5242880},
	}

	// Set the Content-Type header to signal we are sending JSON.
	w.Header().Set("Content-Type", "application/json")

	// Use json.NewEncoder to encode our 'files' slice directly to the
	// ResponseWriter. This is efficient as it avoids creating an
	// intermediate buffer for the JSON data.
	if err := json.NewEncoder(w).Encode(files); err != nil {
		// In a real app, you'd have more robust error handling.
		log.Printf("Error encoding files to JSON: %v", err)
	}
}

// newRouter creates a new Chi router and registers the routes.
func newRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", homeHandler)

	// We create a new group for our API routes under the "/api/v1" prefix.
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/files", listFilesHandler)
		// Future API routes like POST /files will go here.
	})

	return r
}

func main() {
	r := newRouter()

	port := "8080"

	// Log that the server is starting.
	log.Printf("🚀 Server listening on http://localhost:%s", port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}