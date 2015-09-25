package main

import (
	"net/http"

	"github.com/go-zoo/bone"
)

func main() {
	muxx := bone.New()
	sub := bone.New()
	sub.GetFunc("/user/:id", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(bone.GetValue(req, "id")))
		return
	})
	muxx.Get("/api", sub)

	http.ListenAndServe(":8080", muxx)
}
