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

// Route content the required information for a valid route
type Route struct {
	Path    string
	Size    int
	pattern pattern
	handler http.Handler
	Method  string
}

// pattern content the required information for the route pattern
type pattern struct {
	Exist bool
	Id    string
	Pos   int
}

// NewRoute return a pointer to a Route instance and call save() on it
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

// Check if the request match the route pattern
func (r *Route) Matcher(path string) (url.Values, bool) {
	ss := strings.Split(path, "/")

	if len(ss) == r.Size && ss[r.Size-1] != "" {
		uV := url.Values{}
		uV.Add(r.pattern.Id, ss[r.pattern.Pos])
		return uV, true
	}

	return nil, false
}

// Set the route method to Get
func (r *Route) Get() *Route {
	r.Method = "GET"
	return r
}

// Set the route method to Post
func (r *Route) Post() *Route {
	r.Method = "POST"
	return r
}

// Set the route method to Put
func (r *Route) Put() *Route {
	r.Method = "PUT"
	return r
}

// Set the route method to Delete
func (r *Route) Delete() *Route {
	r.Method = "DELETE"
	return r
}

// Set the route method to Head
func (r *Route) Head() *Route {
	r.Method = "HEAD"
	return r
}

// Set the route method to Patch
func (r *Route) Patch() *Route {
	r.Method = "PATCH"
	return r
}

// Set the route method to Options
func (r *Route) Options() *Route {
	r.Method = "OPTIONS"
	return r
}

// Only using this in squiidz/fur package
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

// Check if the request respect the route method if provided.
func (r *Route) MethCheck(req *http.Request) bool {
	if r.Method != "" {
		if req.Method == r.Method {
			return true
		} else {
			return false
		}
	} else {
		return true
	}
}
