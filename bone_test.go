package bone

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test if the route is valid
func TestRouting(t *testing.T) {
	mux := New()
	call := false
	mux.Get("/a/:id", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))

	r, _ := http.NewRequest("GET", "/b/123", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if call {
		t.Error("handler should not be called")
	}
}

// Test if the http method is valid
func TestRoutingMethod(t *testing.T) {
	mux := New()
	call := false
	mux.Get("/t", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))

	r, _ := http.NewRequest("POST", "/t", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if call {
		t.Error("response to a wrong method")
	}
}

// Test if the mux don't handle by prefix
func TestRoutingPath(t *testing.T) {
	mux := New()
	call := false
	mux.Get("/t", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))
	mux.Get("/t/x", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))

	r, _ := http.NewRequest("GET", "/t/x", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if call {
		t.Error("response with the wrong path")
	}
}

func TestRoutingVariable(t *testing.T) {
	var (
		expected = "variable"
		got      string
		mux      = New()
		w        = httptest.NewRecorder()
	)
	mux.Get("/:var", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.URL.Query().Get("var")
	}))

	r, err := http.NewRequest("GET", fmt.Sprintf("/%s", expected), nil)
	if err != nil {
		t.Fatal(err)
	}
	mux.ServeHTTP(w, r)

	if got != expected {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}
