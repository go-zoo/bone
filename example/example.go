package main

import (
	"log"
	"net/http"

	"github.com/squiidz/bone"
)

func main() {
	// New mux instance
	mux := bone.New()
	// Custom 404
	mux.NotFound(Handler404)
	// Handle with any http method, Handle takes http.Handler as argument.
	mux.Handle("/index", homeHandler)
	// Get, Post etc... takes http.HandlerFunc as argument.
	mux.Post("/home", homeHandler)

	// Start Listening
	http.ListenAndServe(":8080", mux)
}

func homeHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("WELCOME HOME"))
}

func Handler404(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("These are not the droids you're looking for ..."))
}
