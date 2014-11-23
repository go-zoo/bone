package bone

import (
	"net/http"
)

type Mux struct {
	handlers map[string]map[string]http.Handler
	NotFound http.Handler
}

func NewMux() *Mux {
	return &Mux{make(map[string]map[string]http.Handler), nil}
}

func (m *Mux) SetNotFound(h http.HandlerFunc) {
	m.NotFound = h
}

func (m *Mux) handle(meth string, path string, h http.Handler) {
	if path != "" {
		if m.handlers[path] != nil {
			m.handlers[path][meth] = h
		} else {
			m.handlers[path] = make(map[string]http.Handler)
			m.handlers[path][meth] = h
		}
	} else {
		panic("Non-Valid Path")
	}
}

// GET set Handler valid method to GET only.
func (m *Mux) GET(path string, h http.Handler) {
	m.handle("GET", path, h)
}

// POST set Handler valid method to POST only.
func (m *Mux) POST(path string, h http.Handler) {
	m.handle("POST", path, h)
}

// DELETE set Handler valid method to DELETE only.
func (m *Mux) DELETE(path string, h http.Handler) {
	m.handle("DELETE", path, h)
}

// PUT set Handler valid method to PUT only.
func (m *Mux) PUT(path string, h http.Handler) {
	m.handle("PUT", path, h)
}

func (m *Mux) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var path = req.URL.Path
	var meth = req.Method

	if h, ok := m.handlers[path]; ok {
		if m, ok := h[meth]; ok {
			m.ServeHTTP(rw, req)
		} else {
			http.Error(rw, "Bad HTTP Method", http.StatusBadRequest)
		}
	} else {
		switch m.NotFound {
		case nil:
			http.NotFound(rw, req)
		default:
			m.NotFound.ServeHTTP(rw, req)
		}
	}

}
