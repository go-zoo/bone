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

// Test the custom not handler handler sets 404 error code
func TestNotFoundCustomHandlerSends404(t *testing.T) {
	mux := New()
	mux.NotFound(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("These are not the droids you're looking for ..."))
	})

	r, _ := http.NewRequest("GET", "/b/123", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Errorf("expecting error code 404, got %v", w.Code)
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
		got = GetValue(r, "var")
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

func TestRoutingVerbs(t *testing.T) {
	var (
		methods = []string{"DELETE", "GET", "HEAD", "PUT", "POST", "PATCH", "OPTIONS", "HEAD"}
		path    = "/"
		h       = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	)
	for _, meth := range methods {
		m := New()
		switch meth {
		case "DELETE":
			m.Delete(path, h)
		case "GET":
			m.Get(path, h)
		case "HEAD":
			m.Head(path, h)
		case "POST":
			m.Post(path, h)
		case "PUT":
			m.Put(path, h)
		case "PATCH":
			m.Patch(path, h)
		case "OPTIONS":
			m.Options(path, h)
		}
		s := httptest.NewServer(m)
		req, err := http.NewRequest(meth, s.URL, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Fatalf("%s: HTTP %d", meth, resp.StatusCode)
		}
		s.Close()
	}
}

func TestRoutingSlash(t *testing.T) {
	mux := New()
	call := false
	mux.Get("/", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !call {
		t.Error("root not serve")
	}
}
