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
)

// If no pattern are set in the route.
// Handle the bad request
func (m *Mux) BadRequest(rw http.ResponseWriter, req *http.Request) {
	if m.notFound != nil {
		rw.WriteHeader(http.StatusNotFound)
		m.notFound(rw, req)
	} else {
		http.NotFound(rw, req)
	}
}

// Check if the path don't end with a /
func valid(path string) bool {
	if len(path) > 1 && path[len(path)-1:] == "/" {
		return false
	}
	return true
}

// Check if the requested route is for a static file
func (m *Mux) isStatic(s string) bool {
	sl := len(s)
	if sl > 1 && s[sl-1:] == "/" {
		return true
	}
	return false
}

// Debugging function
func (m *Mux) inspect(meth string) {
	for i, r := range m.Routes[meth] {
		fmt.Printf("#%d => %s\n", i+1, r.Path)
	}
}

// Insert the url value into the vars stack
func (r *Route) insert(req *http.Request, uv url.Values) {
	for k, _ := range uv {
		r.Pattern.Value[k] = uv.Get(k)
	}
	VARS[req] = r
}

// Return the key value, of the current *http.Request
func GetValue(req *http.Request, key string) string {
	return VARS[req].Pattern.Value[key]
}
