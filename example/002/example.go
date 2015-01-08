package main

import (
	"io/ioutil"
	"net/http"

	"github.com/go-zoo/bone"
)

func main() {
	mux := bone.New()

	mux.Get("/", http.HandlerFunc(defaultHandler))
	mux.Get("/test", http.HandlerFunc(defaultHandler))
	mux.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir("assets"))))

	http.ListenAndServe(":8080", mux)
}

func defaultHandler(rw http.ResponseWriter, req *http.Request) {
	file, _ := ioutil.ReadFile("index.html")
	rw.Write(file)
}
