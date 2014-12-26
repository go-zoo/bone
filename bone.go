/********************************
*** Multiplexer for Go        ***
*** Code is under MIT license ***
*** Code by CodingFerret      ***
*** github.com/squiidz        ***
*********************************/

package bone

import (
	"net/http"
	"sort"
)

// Mux have routes and a notFound handler
// Route: all the registred route
// notFound: 404 handler, default http.NotFound if not provided
type Mux struct {
	Routes   map[string][]*Route
	Static   []*Route
	notFound http.HandlerFunc
}

var (
	METHOD = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "PATCH", "OPTIONS"}
	VARS   = make(map[*http.Request]*Route)
)

// New create a pointer to a Mux instance
func New() *Mux {
	return &Mux{
		Routes: make(map[string][]*Route),
	}
}

// Serve http request
func (m *Mux) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	reqPath := req.URL.Path
	reqLen := len(reqPath)

	// Check if the request path doesn't end with /
	if !valid(reqPath) {
		http.Redirect(rw, req, reqPath[:reqLen-1], http.StatusMovedPermanently)
		return
	}
	// Loop over all the registred route.
	for _, r := range m.Routes[req.Method] {
		// If the route is equal to the request path.
		if reqPath == r.Path {
			r.Handler.ServeHTTP(rw, req)
			return
		} else if r.Pattern.Exist {
			if v, ok := r.Match(req.URL.Path); ok {
				r.insert(req, v)
				r.Handler.ServeHTTP(rw, req)
				return
			}
			// If no pattern are set in the route.
			continue
		}
		continue
	}
	// If no valid Route found, check for static file
	for _, s := range m.Static {
		if reqLen >= s.Size && reqPath[:s.Size] == s.Path {
			s.ServeHTTP(rw, req)
			return
		}
		continue
	}
	m.BadRequest(rw, req)
}

// HandleFunc is use to pass a func(http.ResponseWriter, *Http.Request) instead of http.Handler
func (m *Mux) HandleFunc(path string, handler http.HandlerFunc) {
	m.Handle(path, handler)
}

// Handle add a new route to the Mux without a HTTP method
func (m *Mux) Handle(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	if m.isStatic(path) {
		m.Static = append(m.Static, r.Get())
		return
	} else {
		for _, mt := range METHOD {
			m.Routes[mt] = append(m.Routes[mt], r)
			sort.Sort(ByLength(m.Routes[mt]))
		}
	}
}

// Get add a new route to the Mux with the Get method
func (m *Mux) Get(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes["GET"] = append(m.Routes["GET"], r.Get())
	sort.Sort(ByLength(m.Routes["GET"]))
}

// Post add a new route to the Mux with the Post method
func (m *Mux) Post(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes["POST"] = append(m.Routes["POST"], r.Post())
	sort.Sort(ByLength(m.Routes["POST"]))
}

// Put add a new route to the Mux with the Put method
func (m *Mux) Put(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes["PUT"] = append(m.Routes["PUT"], r.Put())
	sort.Sort(ByLength(m.Routes["PUT"]))
}

// Delete add a new route to the Mux with the Delete method
func (m *Mux) Delete(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes["DELETE"] = append(m.Routes["DELETE"], r.Delete())
	sort.Sort(ByLength(m.Routes["DELETE"]))
}

// Head add a new route to the Mux with the Head method
func (m *Mux) Head(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes["HEAD"] = append(m.Routes["HEAD"], r.Head())
	sort.Sort(ByLength(m.Routes["HEAD"]))
}

// Patch add a new route to the Mux with the Patch method
func (m *Mux) Patch(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes["PATCH"] = append(m.Routes["PATCH"], r.Patch())
	sort.Sort(ByLength(m.Routes["PATCH"]))
}

// Options add a new route to the Mux with the Options method
func (m *Mux) Options(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes["OPTIONS"] = append(m.Routes["OPTIONS"], r.Options())
	sort.Sort(ByLength(m.Routes["OPTIONS"]))
}

// Set the mux custom 404 handler
func (m *Mux) NotFound(handler http.HandlerFunc) {
	m.notFound = handler
}
