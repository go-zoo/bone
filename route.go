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
// Path: is the Route URL
// Size: is the length of the path
// Token: is the value of each part of the path, split by /
// pattern: is content information about the route, if it's have a route variable
// handler: is the handler who handle this route
// Method: define HTTP method on the route
type Route struct {
	Path    string
	Size    int
	Token   token
	pattern pattern
	handler http.Handler
	Method  string
}

// token content all value of a spliting route path
// tokens: string value of each token
// size: number of token
type token struct {
	tokens []string
	size   int
}

type ByLength []*Route

func (b ByLength) Len() int {
	return len(b)
}

func (b ByLength) Swap(i int, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByLength) Less(i int, j int) bool {
	return b[i].Size < b[j].Size
}

// pattern content the required information for the route pattern
// Exist: check if a variable was declare on the route
// Id: the name of the variable
// Pos: postition of var in the route path
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
	r.Token.tokens = strings.Split(r.Path, "/")
	for i, s := range r.Token.tokens {
		if len(s) >= 1 {
			if s[:1] == ":" {
				r.pattern.Exist = true
				r.pattern.Id = s[1:]
				r.pattern.Pos = i
			}
		}
	}
	r.Size = len(r.Path)
	r.Token.size = len(r.Token.tokens)
}

// Info is only used for debugging
func (r *Route) Info() {
	fmt.Printf("Path :         %s\n", r.Path)
	fmt.Printf("Size : 		   %d\n", r.Size)
	fmt.Printf("Have Pattern : %t\n", r.pattern.Exist)
	fmt.Printf("ID :           %s\n", r.pattern.Id)
	fmt.Printf("Position :     %d\n", r.pattern.Pos)
	fmt.Printf("Method :       %s\n", r.Method)
}

// Check if the request match the route pattern
func (r *Route) Match(path string) (url.Values, bool) {
	ss := strings.Split(path, "/")
	if len(ss) == r.Token.size && ss[r.Token.size-1] != "" {
		if r.Path[:r.pattern.Pos] == path[:r.pattern.Pos] {
			uV := url.Values{}
			uV.Add(r.pattern.Id, ss[r.pattern.Pos])
			return uV, true
		}
	}
	return nil, false
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
