/********************************
*** Multiplexer for Go        ***
*** Code is under MIT license ***
*** Code by CodingFerret      ***
*** github.com/squiidz        ***
*********************************/

package bone

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Route struct {
	Path    string
	Size    int
	pattern pattern
	handler http.Handler
	Method  string
}

type pattern struct {
	Exist bool
	Id    string
	Pos   int
}

func NewRoute(url string, h http.Handler) *Route {
	r := &Route{Path: url, handler: h}
	r.save()
	return r
}

// Save, set automaticly the the Route.Size and Route.pattern value
func (r *Route) save() {
	subs := strings.Split(r.Path, "/")
	for i, s := range subs {
		if len(s) >= 1 {
			if s[:1] == ":" {
				r.pattern.Exist = true
				r.pattern.Id = s[1:]
				r.pattern.Pos = i
			}
		}
	}
	r.Size = len(subs)
	return

}

// Info is only use for debugging
func (r *Route) Info() {
	fmt.Printf("Path : %s\n", r.Path)
	fmt.Printf("Size : %d\n", r.Size)
	fmt.Printf("Have Pattern : %t\n", r.pattern.Exist)
	fmt.Printf("ID : %s\n", r.pattern.Id)
	fmt.Printf("Position : %d\n", r.pattern.Pos)
	fmt.Printf("Method : %s\n", r.Method)
}

func (r *Route) Matcher(path string) (url.Values, bool) {
	ss := strings.Split(path, "/")

	if len(ss) == r.Size && ss[r.Size-1] != "" {
		uV := url.Values{}
		uV.Add(r.pattern.Id, ss[r.pattern.Pos])
		return uV, true
	}

	return nil, false
}

func (r *Route) Get() {
	r.Method = "GET"
}

func (r *Route) Post() {
	r.Method = "POST"
}

func (r *Route) Put() {
	r.Method = "PUT"
}

func (r *Route) Delete() {
	r.Method = "DELETE"
}

func (r *Route) Head() {
	r.Method = "HEAD"
}

func (r *Route) Patch() {
	r.Method = "PATCH"
}

func (r *Route) Options() {
	r.Method = "OPTIONS"
}

func (r Route) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	if r.Method != "" {

		if req.Method == r.Method {
			r.handler.ServeHTTP(rw, req)
		} else {
			http.NotFound(rw, req)
		}

	} else {
		r.handler.ServeHTTP(rw, req)
	}

	// DEBUG r.Info()
}
