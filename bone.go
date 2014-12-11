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

func NewMux() *Mux {
	return &Mux{}
}

func (m *Mux) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	reqPath := req.URL.Path
	if !valid(reqPath) {
		http.Redirect(rw, req, reqPath[:len(reqPath)-1], http.StatusMovedPermanently)
		return
	}

	for _, r := range m.Routes {
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

	if m.notFound != nil {
		m.notFound(rw, req)
	} else {
		http.NotFound(rw, req)
	}

}

func valid(path string) bool {
	pathLen := len(path)

	if pathLen > 1 && path[pathLen-1:] == "/" {
		return false
	}
	return true
}

func (m *Mux) Handle(s string, h http.Handler) {
	m.Routes = append(m.Routes, NewRoute(s, h))
}

func (m *Mux) NotFound(h http.HandlerFunc) {
	m.notFound = h
}
