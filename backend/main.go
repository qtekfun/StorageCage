// backend/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
)

// homeHandler handles requests to the root ("/") path.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "StorageCage API server is running! 👋")
}

// newRouter creates a new Chi router and registers the routes.
func newRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", homeHandler)
	return r
}

func main() {
	r := newRouter()

	port := "8080"

	// Log that the server is starting.
	log.Printf("🚀 Server listening on http://localhost:%s", port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}