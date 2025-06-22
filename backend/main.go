// backend/main.go (versión actualizada)
package main

import (
	"log"
	"net/http"
)

func main() {
	// Create a new config instance.
	// For now, we hardcode the storage path. Later, this could come from a file or env vars.
	config := AppConfig{
		StorageDir: "./data", // We'll store files in a 'data' directory.
	}

	router := newRouter(&config) // Pass the config to the router
	port := "8080"

	log.Printf("🚀 Server listening on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
