/********************************
*** Multiplexer for Go        ***
*** Bone is under MIT license ***
*** Code by CodingFerret      ***
*** github.com/squiidz        ***
*********************************/

package bone

import (
	"fmt"
	"net/http"
	"net/url"
)

// BadRequest handle every wring request.
func (m *Mux) BadRequest(rw http.ResponseWriter, req *http.Request) {
	if m.notFound != nil {
		rw.WriteHeader(http.StatusNotFound)
		m.notFound(rw, req)
	} else {
		http.NotFound(rw, req)
	}
}

// Check if the path don't end with a /
func (m *Mux) valid(path string) bool {
	plen := len(path)
	if plen > 1 && path[plen-1:] == "/" && !m.inStatic(path) {
		return false
	}
	return true
}

func (m *Mux) inStatic(s string) bool {
	for k := range m.Static {
		if k == s {
			return true
		}
	}
	return false
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

// info is only used for debugging
func (r *Route) info() {
	fmt.Printf("Path :         %s\n", r.Path)
	fmt.Printf("Size : 		   %d\n", r.Size)
	fmt.Printf("Have Pattern : %t\n", r.Pattern.Exist)
	fmt.Printf("ID :           %s\n", r.Pattern.ID)
	fmt.Printf("Position :     %d\n", r.Pattern.Pos)
	fmt.Printf("Method :       %s\n", r.Method)
}

// Insert the url value into the vars stack
func (r *Route) insert(req *http.Request, uv url.Values) {
	for k := range uv {
		r.Pattern.Value[k] = uv.Get(k)
	}
	vars[req] = r
}

// GetValue Return the key value, of the current *http.Request
func GetValue(req *http.Request, key string) string {
	return vars[req].Pattern.Value[key]
}
