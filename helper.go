/********************************
*** Multiplexer for Go        ***
*** Bone is under MIT license ***
*** Code by CodingFerret      ***
*** github.com/go-zoo         ***
*********************************/

package bone

import "net/http"

func (m *Mux) parse(rw http.ResponseWriter, req *http.Request) bool {
	for _, r := range m.Routes[req.Method] {
		if r.Atts != 0 {
			if r.Atts&SUB != 0 {
				if len(req.URL.Path) >= r.Size {
					if req.URL.Path[:r.Size] == r.Path {
						req.URL.Path = req.URL.Path[r.Size:]
						r.Handler.ServeHTTP(rw, req)
						return true
					}
				}
			}
			if r.Match(req) {
				r.Handler.ServeHTTP(rw, req)
				vars.Lock()
				delete(vars.v, req)
				vars.Unlock()
				return true
			}
		}
		if req.URL.Path == r.Path {
			r.Handler.ServeHTTP(rw, req)
			return true
		}
	}
	return false
}

// StaticRoute check if the request path is for Static route
func (m *Mux) staticRoute(rw http.ResponseWriter, req *http.Request) bool {
	for _, s := range m.Routes[static] {
		if len(req.URL.Path) >= s.Size {
			if req.URL.Path[:s.Size] == s.Path {
				s.Handler.ServeHTTP(rw, req)
				return true
			}
		}
	}
	return false
}

// HandleNotFound handle when a request does not match a registered handler.
func (m *Mux) HandleNotFound(rw http.ResponseWriter, req *http.Request) {
	if m.notFound != nil {
		m.notFound.ServeHTTP(rw, req)
	} else {
		http.NotFound(rw, req)
	}
}

// Check if the path don't end with a /
func (m *Mux) validate(rw http.ResponseWriter, req *http.Request) bool {
	plen := len(req.URL.Path)
	if plen > 1 && req.URL.Path[plen-1:] == "/" {
		cleanURL(&req.URL.Path)
		rw.Header().Set("Location", req.URL.Path)
		rw.WriteHeader(http.StatusFound)
	}
	// Retry to find a route that match
	return m.parse(rw, req)
}

func valid(path string) bool {
	plen := len(path)
	if plen > 1 && path[plen-1:] == "/" {
		return false
	}
	return true
}

// Clean url path
func cleanURL(url *string) {
	ulen := len((*url))
	if ulen > 1 {
		if (*url)[ulen-1:] == "/" {
			*url = (*url)[:ulen-1]
			cleanURL(url)
		}
	}
}

// GetValue return the key value, of the current *http.Request
func GetValue(req *http.Request, key string) string {
	vars.RLock()
	value := vars.v[req][key]
	vars.RUnlock()
	return value
}

// GetAllValues return the req PARAMs
func GetAllValues(req *http.Request) map[string]string {
	vars.RLock()
	values := vars.v[req]
	vars.RUnlock()
	return values
}

// This function returns the route of given Request
func (m *Mux) GetRequestRoute(req *http.Request) string {
	cleanURL(&req.URL.Path)
	for _, r := range m.Routes[req.Method] {
		if r.Atts != 0 {
			if r.Atts&SUB != 0 {
				if len(req.URL.Path) >= r.Size {
					if req.URL.Path[:r.Size] == r.Path {
						return r.Path
					}
				}
			}
			if r.Match(req) {
				return r.Path
			}
		}
		if req.URL.Path == r.Path {
			return r.Path
		}
	}

	for _, s := range m.Routes[static] {
		if len(req.URL.Path) >= s.Size {
			if req.URL.Path[:s.Size] == s.Path {
				return s.Path
			}
		}
	}

	return "404-NotFound"
}
