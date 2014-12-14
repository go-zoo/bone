/********************************
*** Multiplexer for Go        ***
*** Code is under MIT license ***
*** Code by CodingFerret      ***
*** github.com/squiidz        ***
*********************************/

package bone

import (
	"net/http"
	"strings"
)

// Mux have routes and a notFound handler
// Route: all the registred route
// notFound: 404 handler, default http.NotFound if not provided
type Mux struct {
	Routes   []*Route
	notFound http.HandlerFunc
}

// New create a pointer to a Mux instance
func New() *Mux {
	return &Mux{}
}

// Serve http request
func (m *Mux) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	reqPath := req.URL.Path
	// Check if the request path doesn't end with /
	if !valid(reqPath) {
		http.Redirect(rw, req, reqPath[:len(reqPath)-1], http.StatusMovedPermanently)
		return
	}
	// Loop over all the registred route.
	for _, r := range m.Routes {
		// Check if the request method is valid.
		if !r.MethCheck(req) {
			continue
		}
		// If the route have a pattern.
		if r.pattern.Exist {
			if v, ok := r.Match(req.URL.Path); ok {
				req.URL.RawQuery = v.Encode() + "&" + req.URL.RawQuery
				r.handler.ServeHTTP(rw, req)
				return
			}
			continue
			// If no pattern are set in the route.
		} else {
			if len(req.URL.Path) == r.Size && req.URL.Path[:r.Size] == r.Path {
				r.handler.ServeHTTP(rw, req)
				return
			} else if fileExt(req.URL.Path) {
				r.handler.ServeHTTP(rw, req)
				return
			}
			continue
		}
	}
	m.BadRequest(rw, req)
}

// If no pattern are set in the route.
// Handle the bad request
func (m *Mux) BadRequest(rw http.ResponseWriter, req *http.Request) {
	if m.notFound != nil {
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
func fileExt(s string) bool {
	parts := strings.Split(s, "/")
	if strings.Contains(parts[len(parts)-1], ".") {
		return true
	}
	return false
}

// Handle add a new route to the Mux without a HTTP method
func (m *Mux) Handle(path string, handler http.Handler) {
	m.Routes = append(m.Routes, NewRoute(path, handler))
}

// Get add a new route to the Mux with the Get method
func (m *Mux) Get(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes = append(m.Routes, r.Get())
}

// Post add a new route to the Mux with the Post method
func (m *Mux) Post(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes = append(m.Routes, r.Post())
}

// Put add a new route to the Mux with the Put method
func (m *Mux) Put(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes = append(m.Routes, r.Put())
}

// Delete add a new route to the Mux with the Delete method
func (m *Mux) Delete(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes = append(m.Routes, r.Delete())
}

// Head add a new route to the Mux with the Head method
func (m *Mux) Head(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes = append(m.Routes, r.Head())
}

// Patch add a new route to the Mux with the Patch method
func (m *Mux) Patch(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes = append(m.Routes, r.Patch())
}

// Options add a new route to the Mux with the Options method
func (m *Mux) Options(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes = append(m.Routes, r.Options())
}

// Set the mux custom 404 handler
func (m *Mux) NotFound(handler http.HandlerFunc) {
	m.notFound = handler
}
