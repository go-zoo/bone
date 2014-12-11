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

type Mux struct {
	Routes   []*Route
	notFound http.HandlerFunc
}

func New() *Mux {
	return &Mux{}
}

func (m *Mux) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	reqPath := req.URL.Path
	if !valid(reqPath) {
		http.Redirect(rw, req, reqPath[:len(reqPath)-1], http.StatusMovedPermanently)
		return
	}

	for _, r := range m.Routes {

		if !r.MethCheck(req) {
			m.BadRequest(rw, req)
			return
		}

		if r.pattern.Exist {
			if v, ok := r.Matcher(req.URL.Path); ok {
				req.URL.RawQuery = v.Encode() + "&" + req.URL.RawQuery
				r.handler.ServeHTTP(rw, req)
				return
			}
			continue
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

func (m *Mux) BadRequest(rw http.ResponseWriter, req *http.Request) {
	if m.notFound != nil {
		m.notFound(rw, req)
	} else {
		http.NotFound(rw, req)
	}
}

func valid(path string) bool {
	if len(path) > 1 && path[len(path)-1:] == "/" {
		return false
	}
	return true
}

func (m *Mux) Handle(s string, h http.Handler) {
	m.Routes = append(m.Routes, NewRoute(s, h))
}

func (m *Mux) Get(s string, h http.HandlerFunc) {
	r := NewRoute(s, h)
	r.Get()
	m.Routes = append(m.Routes, r)
}

func (m *Mux) Post(s string, h http.HandlerFunc) {
	r := NewRoute(s, h)
	r.Post()
	m.Routes = append(m.Routes, r)
}

func (m *Mux) Put(s string, h http.HandlerFunc) {
	r := NewRoute(s, h)
	r.Put()
	m.Routes = append(m.Routes, r)
}

func (m *Mux) Delete(s string, h http.HandlerFunc) {
	r := NewRoute(s, h)
	r.Delete()
	m.Routes = append(m.Routes, r)
}

func (m *Mux) Head(s string, h http.HandlerFunc) {
	r := NewRoute(s, h)
	r.Head()
	m.Routes = append(m.Routes, r)
}

func (m *Mux) Patch(s string, h http.HandlerFunc) {
	r := NewRoute(s, h)
	r.Patch()
	m.Routes = append(m.Routes, r)
}

func (m *Mux) Options(s string, h http.HandlerFunc) {
	r := NewRoute(s, h)
	r.Options()
	m.Routes = append(m.Routes, r)
}

func (m *Mux) NotFound(h http.HandlerFunc) {
	m.notFound = h
}
