/********************************
*** Multiplexer for Go        ***
*** Bone is under MIT license ***
*** Code by CodingFerret      ***
*** github.com/go-zoo         ***
*********************************/

package bone

import "net/http"

// Mux have routes and a notFound handler
// Route: all the registred route
// notFound: 404 handler, default http.NotFound if not provided
type Mux struct {
	Routes   map[string][]*Route
	Static   map[string]*Route
	notFound http.HandlerFunc
}

var (
	method = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "PATCH", "OPTIONS"}
	vars   = make(map[*http.Request]*Route)
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

	// Check if the request path doesn't end with /
	if !m.valid(req.URL.Path) {
		if key, ok := m.isStatic(req.URL.Path); ok {
			m.Static[key].Handler.ServeHTTP(rw, req)
			return
		}
		for !m.valid(req.URL.Path) {
			req.URL.Path = req.URL.Path[:len(req.URL.Path)-1]
		}

		rw.Header().Set("Location", req.URL.Path)
		rw.WriteHeader(http.StatusFound)
	}

	// Loop over all the registred route.
	for _, r := range m.Routes[req.Method] {
		// If the route is equal to the request path.
		if req.URL.Path == r.Path && !r.Pattern.Exist {
			r.Handler.ServeHTTP(rw, req)
			return
		} else if r.Pattern.Exist {
			if v, ok := r.Match(req.URL.Path); ok {
				r.insert(req, v)
				r.Handler.ServeHTTP(rw, req)
				return
			}
			continue
		}
		continue
	}
	// If no valid Route found, check for static file
	if key, ok := m.isStatic(req.URL.Path); ok {
		m.Static[key].Handler.ServeHTTP(rw, req)
		return
	}
	m.HandleNotFound(rw, req)
}

// Handle add a new route to the Mux without a HTTP method
func (m *Mux) Handle(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	if !m.valid(path) {
		m.Static[path] = r
		return
	}
	for _, mt := range method {
		m.register(mt, path, handler)
	}
}

// HandleFunc is use to pass a func(http.ResponseWriter, *Http.Request) instead of http.Handler
func (m *Mux) HandleFunc(path string, handler http.HandlerFunc) {
	m.Handle(path, handler)
}

// Get add a new route to the Mux with the Get method
func (m *Mux) Get(path string, handler http.Handler) {
	m.register("GET", path, handler)
}

// Post add a new route to the Mux with the Post method
func (m *Mux) Post(path string, handler http.Handler) {
	m.register("POST", path, handler)
}

// Put add a new route to the Mux with the Put method
func (m *Mux) Put(path string, handler http.Handler) {
	m.register("PUT", path, handler)
}

// Delete add a new route to the Mux with the Delete method
func (m *Mux) Delete(path string, handler http.Handler) {
	m.register("DELETE", path, handler)
}

// Head add a new route to the Mux with the Head method
func (m *Mux) Head(path string, handler http.Handler) {
	m.register("HEAD", path, handler)
}

// Patch add a new route to the Mux with the Patch method
func (m *Mux) Patch(path string, handler http.Handler) {
	m.register("PATCH", path, handler)
}

// Options add a new route to the Mux with the Options method
func (m *Mux) Options(path string, handler http.Handler) {
	m.register("OPTIONS", path, handler)
}

// NotFound the mux custom 404 handler
func (m *Mux) NotFound(handler http.HandlerFunc) {
	m.notFound = handler
}

// Register the new route in the router with the provided method and handler
func (m *Mux) register(method string, path string, handler http.Handler) {
	if m.valid(path) {
		r := NewRoute(path, handler)
		m.Routes[method] = append(m.Routes[method], r)
		byLength(m.Routes[method]).Sort()
		return
	}
	m.Handle(path, handler)
}
