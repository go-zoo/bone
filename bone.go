package bone

import "net/http"

type Mux struct {
	Handlers map[string]http.Handler
	NotFound http.Handler
}

func NewMux() *Mux {
	return &Mux{make(map[string]http.Handler), nil}
}

func (m *Mux) SetNotFound(h http.Handler) {
	m.NotFound = h
}

func (m *Mux) Handle(path string, h http.Handler) {
	if path != "" {
		m.Handlers[path] = h
	} else {
		panic("Non-Valid Path")
	}
}

func (m *Mux) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var valid = ""
	for k, _ := range m.Handlers {
		if req.URL.Path == k {
			valid = k
		}
	}
	if valid != "" {
		m.Handlers[valid].ServeHTTP(rw, req)
	} else {
		if m.NotFound != nil {
			m.NotFound.ServeHTTP(rw, req)
		} else {
			http.NotFound(rw, req)
		}
	}

}
