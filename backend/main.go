// backend/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
)

// homeHandler handles requests to the root ("/") path.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "StorageCage API server is running! 👋")
}

func main() {
	// Register the homeHandler for the "/" route.
	http.HandleFunc("/", homeHandler)

	port := "8080"

	// Log that the server is starting.
	log.Printf("🚀 Server listening on http://localhost:%s", port)

	// Start the server and log a fatal error if it fails.
	log.Fatal(http.ListenAndServe(":"+port, nil))
}