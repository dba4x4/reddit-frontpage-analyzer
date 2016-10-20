package server

import (
	"log"
	"net/http"
)

// Start the web server that serves the API
func Start() {
	log.Println("Listening on http://127.0.0.1:8080")
	router := newRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
