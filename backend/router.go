package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// newRouter sets up our application's routes.
func newRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", homeHandler)

	// API routes group
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/files", listFilesHandler)
	})

	return r
}
