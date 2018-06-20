// +build go1.7

package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-zoo/bone"
)

func main() {
	mux := bone.New()
	mux.CaseSensitive = true
	mux.RegisterValidatorFunc("isNum", func(s string) bool {
		if _, err := strconv.Atoi(s); err == nil {
			return true
		}
		return false
	})

	mux.RegisterValidatorFunc("biggerThan1000", func(s string) bool {
		if num, err := strconv.Atoi(s); err == nil {
			if num >= 1000 {
				return true
			}
		}
		return false
	})

	mux.RegisterValidatorFunc("lessThan8", func(s string) bool {
		if len(s) < 8 {
			return true
		}
		return false
	})

	mux.RegisterValidator("exist", &Exist{
		things: []string{"steve", "john", "fee", "charlotte"},
	})

	mux.GetFunc("/ctx/:age|isNum|biggerThan1000/:name|lessThan8|exist", rootHandler)

	http.ListenAndServe(":8080", mux)
}

func rootHandler(rw http.ResponseWriter, req *http.Request) {
	ctx := context.WithValue(req.Context(), "var", bone.GetValue(req, "var"))
	subHandler(rw, req.WithContext(ctx))
}

func subHandler(rw http.ResponseWriter, req *http.Request) {
	vars := bone.GetAllValues(req)
	age := vars["age"]
	name := vars["name"]
	rw.Write([]byte(age + " " + name))
}

type Exist struct {
	things []string
}

func (e *Exist) Validate(s string) bool {
	for _, thing := range e.things {
		if thing == s {
			return true
		}
	}
	return false
}
