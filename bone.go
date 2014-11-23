/**********************************
***      HTTP Router in Go      ***
***  Code is under MIT license  ***
***    Code by CodingFerret     ***
*** 	github.com/squiidz      ***
***********************************/

package bone

import (
	"net/http"
	"net/url"
	"strings"
)

type Handler struct {
	Path string
	http.HandlerFunc
}

type Mux struct {
	handlers map[string][]*Handler
	NotFound http.Handler
}

func NewMux() *Mux {
	return &Mux{make(map[string][]*Handler), nil}
}

func (m *Mux) SetNotFound(h http.HandlerFunc) {
	m.NotFound = h
}

func (m *Mux) handle(meth string, h *Handler) {
	if h.Path != "" {
		m.handlers[meth] = append(m.handlers[meth], h)
	} else {
		panic("Non-Valid Path")
	}
}

// GET set Handler valid method to GET only.
func (m *Mux) GET(path string, h http.HandlerFunc) {
	m.handle("GET", &Handler{path, h})
}

// POST set Handler valid method to POST only.
func (m *Mux) POST(path string, h http.HandlerFunc) {
	m.handle("POST", &Handler{path, h})
}

// DELETE set Handler valid method to DELETE only.
func (m *Mux) DELETE(path string, h http.HandlerFunc) {
	m.handle("DELETE", &Handler{path, h})
}

// PUT set Handler valid method to PUT only.
func (m *Mux) PUT(path string, h http.HandlerFunc) {
	m.handle("PUT", &Handler{path, h})
}

func (h *Handler) match(path string) (url.Values, bool) {
	urlVal := url.Values{}
	mp := strings.Split(h.Path[1:], "/")
	rp := strings.Split(path[1:], "/")

	if len(rp) != len(mp) {
		return nil, false
	}

	var rfp string

	for id, val := range mp {
		if len(val) > 1 && val[:1] == "#" {
			urlVal.Add(val[1:], rp[id])
			rfp += "/" + rp[id]
			continue
		}
		rfp += "/" + val
	}

	if rfp != path {
		return nil, false
	}
	return urlVal, true
}

func (m *Mux) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var path = req.URL.Path
	var meth = req.Method
	plen := len(path)

	if plen > 1 && req.URL.Path[plen-1:] == "/" {
		http.Redirect(rw, req, req.URL.Path[:plen-1], 301)
		return
	}
	for _, h := range m.handlers[meth] {
		if vars, ok := h.match(path); ok {
			req.URL.RawQuery = vars.Encode() + "&" + req.URL.RawQuery
			h.ServeHTTP(rw, req)
			return
		}
	}
	switch m.NotFound {
	case nil:
		http.NotFound(rw, req)
		return
	default:
		m.NotFound.ServeHTTP(rw, req)
		return
	}
}
