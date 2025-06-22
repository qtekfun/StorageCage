package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func newRouter(config *AppConfig) http.Handler {
	r := chi.NewRouter()

	// Add the CORS middleware
	r.Use(cors.Handler(cors.Options{
		// This is a permissive configuration for local development.
		// In production, you would want to restrict this to your actual frontend's domain.
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any major browser
	}))

	// The handlers are now methods of the AppConfig struct.
	r.Get("/", config.homeHandler)
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/files", config.listFilesHandler)
		r.Post("/files", config.uploadFileHandler)
		// The {fileName} part is a URL parameter that Chi will capture for us.
		r.Delete("/files/{fileName}", config.deleteFileHandler)
	})
	return r
}
