/********************************
*** Multiplexer for Go        ***
*** Bone is under MIT license ***
*** Code by CodingFerret      ***
*** github.com/go-zoo         ***
*********************************/

package bone

import "net/http"

// HandleNotFound handle when a request does not match a registered handler.
func (m *Mux) HandleNotFound(rw http.ResponseWriter, req *http.Request) {
	if m.notFound != nil {
		rw.WriteHeader(http.StatusNotFound)
		m.notFound.ServeHTTP(rw, req)
	} else {
		http.NotFound(rw, req)
	}
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

// Check if the path don't end with a /
func valid(path string) bool {
	plen := len(path)
	if plen > 1 && path[plen-1:] == "/" {
		return false
	}
	return true
}

// Check if the request path is for Static route
func (m *Mux) StaticRoute(rw http.ResponseWriter, req *http.Request) bool {
	p := req.URL.Path
	for _, s := range m.Static {
		if len(p) >= s.Size && p[:s.Size] == s.Path {
			s.ServeHTTP(rw, req)
			return true
		}
		continue
	}
	return false
}

// GetValue return the key value, of the current *http.Request
func GetValue(req *http.Request, key string) string {
	vars.RLock()
	value := vars.m[req][key]
	vars.RUnlock()
	return value
}

// GetAllValues return the req params
func GetAllValues(req *http.Request) map[string]string {
	return vars.m[req]
}
