/********************************
*** Multiplexer for Go        ***
*** Code is under MIT license ***
*** Code by CodingFerret      ***
*** github.com/squiidz        ***
*********************************/

package bone

import (
	"net/http"
)

// Mux have routes and a notFound handler
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
			if v, ok := r.Matcher(req.URL.Path); ok {
				req.URL.RawQuery = v.Encode() + "&" + req.URL.RawQuery
				r.handler.ServeHTTP(rw, req)
				return
			}
			continue
			// If no pattern are set in the route.
		} else {
			if len(req.URL.Path) >= len(r.Path) && req.URL.Path[:len(r.Path)] == r.Path {
				r.handler.ServeHTTP(rw, req)
				return
			}
			continue
		}
	}

	m.BadRequest(rw, req)
}

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

// Handle add a new route to the Mux without a HTTP method
func (m *Mux) Handle(s string, h http.Handler) {
	m.Routes = append(m.Routes, NewRoute(s, h))
}

// Get add a new route to the Mux with the Get method
func (m *Mux) Get(s string, h http.HandlerFunc) {
	r := NewRoute(s, h)
	m.Routes = append(m.Routes, r.Get())
}

// Post add a new route to the Mux with the Post method
func (m *Mux) Post(s string, h http.HandlerFunc) {
	r := NewRoute(s, h)
	r.Post()
	m.Routes = append(m.Routes, r)
}

// Put add a new route to the Mux with the Put method
func (m *Mux) Put(s string, h http.HandlerFunc) {
	r := NewRoute(s, h)
	r.Put()
	m.Routes = append(m.Routes, r)
}

// Delete add a new route to the Mux with the Delete method
func (m *Mux) Delete(s string, h http.HandlerFunc) {
	r := NewRoute(s, h)
	r.Delete()
	m.Routes = append(m.Routes, r)
}

// Head add a new route to the Mux with the Head method
func (m *Mux) Head(s string, h http.HandlerFunc) {
	r := NewRoute(s, h)
	r.Head()
	m.Routes = append(m.Routes, r)
}

// Patch add a new route to the Mux with the Patch method
func (m *Mux) Patch(s string, h http.HandlerFunc) {
	r := NewRoute(s, h)
	r.Patch()
	m.Routes = append(m.Routes, r)
}

// Options add a new route to the Mux with the Options method
func (m *Mux) Options(s string, h http.HandlerFunc) {
	r := NewRoute(s, h)
	r.Options()
	m.Routes = append(m.Routes, r)
}

// Set the mux custom 404 handler
func (m *Mux) NotFound(h http.HandlerFunc) {
	m.notFound = h
}
