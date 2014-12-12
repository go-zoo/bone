package main

import (
	"net/http"

	"github.com/squiidz/bone"
)

func main() {
	mux := bone.New()

	mux.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir("assets"))))

	http.ListenAndServe(":8080", mux)
}
