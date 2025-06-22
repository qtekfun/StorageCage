package main

import (
	"log"
	"net/http"
)

func main() {
	router := newRouter()
	port := "8080"
	log.Printf("🚀 Server listening on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
