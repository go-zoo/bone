/********************************
*** Multiplexer for Go        ***
*** Bone is under MIT license ***
*** Code by CodingFerret      ***
*** github.com/go-zoo         ***
*********************************/

package bone

import (
	"net/http"
	"sync"
)

// Mux have routes and a notFound handler
// Route: all the registred route
// notFound: 404 handler, default http.NotFound if not provided
type Mux struct {
	Routes   map[string][]*Route
	Static   map[string]*Route
	notFound http.Handler
}

var (
	method = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "PATCH", "OPTIONS"}
	vars   = struct {
		sync.RWMutex
		m map[*http.Request]map[string]string
	}{m: make(map[*http.Request]map[string]string)}
)

// New create a pointer to a Mux instance
func New() *Mux {
	return &Mux{
		Routes: make(map[string][]*Route),
		Static: make(map[string]*Route),
	}
}

// Serve http request
func (m *Mux) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if !m.parse(rw, req) {
		// Check if it's a static ressource
		if !m.StaticRoute(rw, req) {
			// Check if the request path doesn't end with /
			if !valid(req.URL.Path) {
				cleanURL(&req.URL.Path)
				rw.Header().Set("Location", req.URL.Path)
				rw.WriteHeader(http.StatusFound)
				if m.parse(rw, req) {
					return
				}
			}
			m.HandleNotFound(rw, req)
		}
	}
	return
}

func (m *Mux) parse(rw http.ResponseWriter, req *http.Request) bool {
	for _, r := range m.Routes[req.Method] {
		// If the route is equal to the request path.
		if req.URL.Path == r.Path && !r.Params {
			r.Handler.ServeHTTP(rw, req)
			return true
		} else if r.Sub {
			if len(req.URL.Path) >= len(r.Path) {
				if req.URL.Path[:len(r.Path)] == r.Path {
					req.URL.Path = req.URL.Path[len(r.Path):]
					r.Handler.ServeHTTP(rw, req)
					return true
				}
			}
		} else if r.Spc {
			if r.Match(req) {
				r.Handler.ServeHTTP(rw, req)
				vars.Lock()
				delete(vars.m, req)
				vars.Unlock()
				return true
			}
		}
	}
	return false
}

// Handle add a new route to the Mux without a HTTP method
func (m *Mux) Handle(path string, handler http.Handler) {
	for _, mt := range method {
		m.register(mt, path, handler)
	}
}

// HandleFunc is use to pass a func(http.ResponseWriter, *Http.Request) instead of http.Handler
func (m *Mux) HandleFunc(path string, handler http.HandlerFunc) {
	m.Handle(path, handler)
}

// Get add a new route to the Mux with the Get method
func (m *Mux) Get(path string, handler http.Handler) *Route {
	return m.register("GET", path, handler)
}

// Post add a new route to the Mux with the Post method
func (m *Mux) Post(path string, handler http.Handler) *Route {
	return m.register("POST", path, handler)
}

// Put add a new route to the Mux with the Put method
func (m *Mux) Put(path string, handler http.Handler) *Route {
	return m.register("PUT", path, handler)
}

// Delete add a new route to the Mux with the Delete method
func (m *Mux) Delete(path string, handler http.Handler) *Route {
	return m.register("DELETE", path, handler)
}

// Head add a new route to the Mux with the Head method
func (m *Mux) Head(path string, handler http.Handler) *Route {
	return m.register("HEAD", path, handler)
}

// Patch add a new route to the Mux with the Patch method
func (m *Mux) Patch(path string, handler http.Handler) *Route {
	return m.register("PATCH", path, handler)
}

// Options add a new route to the Mux with the Options method
func (m *Mux) Options(path string, handler http.Handler) *Route {
	return m.register("OPTIONS", path, handler)
}

// NotFound the mux custom 404 handler
func (m *Mux) NotFound(handler http.Handler) {
	m.notFound = handler
}

// Register the new route in the router with the provided method and handler
func (m *Mux) register(method string, path string, handler http.Handler) *Route {
	r := NewRoute(path, handler)
	if valid(path) {
		m.Routes[method] = append(m.Routes[method], r)
		return r
	}
	m.Static[path] = r
	return r
}

// SubRoute register a third party router as a subRouter of bone
func (m *Mux) SubRoute(path string, router Router) *Route {
	r := NewRoute(path, router)
	if valid(path) {
		r.Sub = true
		for _, mt := range method {
			m.Routes[mt] = append(m.Routes[mt], r)
		}
		return r
	}
	return nil
}
