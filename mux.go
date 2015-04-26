/********************************
*** Multiplexer for Go        ***
*** Bone is under MIT license ***
*** Code by CodingFerret      ***
*** github.com/go-zoo         ***
*********************************/

package bone

import (
	"net/http"
)

// Get add a new route to the Mux with the Get method
func (m *Mux) GetFunc(path string, handler http.HandlerFunc) {
	m.register("GET", path, handler)
}

// Post add a new route to the Mux with the Post method
func (m *Mux) PostFunc(path string, handler http.HandlerFunc) {
	m.register("POST", path, handler)
}

// Put add a new route to the Mux with the Put method
func (m *Mux) PutFunc(path string, handler http.HandlerFunc) {
	m.register("PUT", path, handler)
}

// Delete add a new route to the Mux with the Delete method
func (m *Mux) DeleteFunc(path string, handler http.HandlerFunc) {
	m.register("DELETE", path, handler)
}

// Head add a new route to the Mux with the Head method
func (m *Mux) HeadFunc(path string, handler http.HandlerFunc) {
	m.register("HEAD", path, handler)
}

// Patch add a new route to the Mux with the Patch method
func (m *Mux) PatchFunc(path string, handler http.HandlerFunc) {
	m.register("PATCH", path, handler)
}

// Options add a new route to the Mux with the Options method
func (m *Mux) OptionsFunc(path string, handler http.HandlerFunc) {
	m.register("OPTIONS", path, handler)
}

// NotFound the mux custom 404 handler
func (m *Mux) NotFoundFunc(handler http.HandlerFunc) {
	m.notFound = handler
}
