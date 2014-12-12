package bone

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouting(t *testing.T) {
	mux := New()
	call := false
	mux.Get("/a/:id", func(http.ResponseWriter, *http.Request) {
		call = true
	})

	r, _ := http.NewRequest("GET", "/b/123", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if call {
		t.Error("handler should not be called")
	}
}
