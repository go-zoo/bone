package main

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
)

func main() {
	boneSub := bone.New()
	gorrilaSub := mux.NewRouter()
	httprouterSub := httprouter.New()

	boneSub.GetFunc("/user/:id", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(bone.GetValue(req, "id")))
		return
	})

	boneSub.GetFunc("/test", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("From Test sub route"))
	})

	boneSub.GetFunc("/test/:id", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("From Test :" + bone.GetValue(req, "id")))
	})

	gorrilaSub.HandleFunc("/gorilla", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("Hello from gorilla mux"))
	})

	httprouterSub.GET("/test", func(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		rw.Write([]byte("Hello from httprouter !"))
	})

	muxx := bone.New()

	muxx.SubRoute("/api", boneSub)
	muxx.SubRoute("/guest", gorrilaSub)
	muxx.SubRoute("/http", httprouterSub)

	http.ListenAndServe(":8080", muxx)
}
