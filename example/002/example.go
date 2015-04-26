package main

import (
	"io/ioutil"
	"net/http"

	"github.com/go-zoo/bone"
)

func main() {
	mux := bone.New()

	mux.Get("/", http.HandlerFunc(defaultHandler))
	mux.Get("/reg/#var^[a-z]$/#var2^[0-9]$", http.HandlerFunc(ShowVar))
	mux.Get("/test", http.HandlerFunc(defaultHandler))
	mux.Get("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir("assets"))))

	http.ListenAndServe(":8080", mux)
}

func defaultHandler(rw http.ResponseWriter, req *http.Request) {
	file, _ := ioutil.ReadFile("index.html")
	rw.Write(file)
}

func ShowVar(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte(bone.GetValue(req, "var") + " - " + bone.GetValue(req, "var2")))
}
