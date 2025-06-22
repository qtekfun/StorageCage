package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func newRouter(config *AppConfig) http.Handler {
	r := chi.NewRouter()
	// The handlers are now methods of the AppConfig struct.
	r.Get("/", config.homeHandler)
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/files", config.listFilesHandler)
	})
	return r
}
