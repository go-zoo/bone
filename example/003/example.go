package main

import (
	"fmt"
	"net/http"

	"github.com/go-zoo/bone"
)

var (
	router = bone.New()
	muxx   = bone.New()
)

func main() {
	muxx.GetFunc("*/test", TestHandler)

	muxx.GetFunc("*/main", TestHandler)

	router.Handle("/index/*", muxx)

	http.ListenAndServe(":8080", router)
}

func TestHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Println(req.RequestURI)
	rw.Write([]byte(req.RequestURI))
}
