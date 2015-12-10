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
	mux.NotFoundFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
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

func TestMultipleRoutingVariables(t *testing.T) {
	var (
		expected1 = "variable1"
		expected2 = "variable2"
		got1      string
		got2      string
		mux       = New()
		w         = httptest.NewRecorder()
	)
	mux.Get("/test/:var1/:var2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got1 = GetValue(r, "var1")
		got2 = GetValue(r, "var2")
	}))

	r, err := http.NewRequest("GET", fmt.Sprintf("/test/%s/%s", expected1, expected2), nil)
	if err != nil {
		t.Fatal(err)
	}
	mux.ServeHTTP(w, r)

	if got1 != expected1 {
		t.Fatalf("expected %s, got %s", expected1, got1)
	}

	if got2 != expected2 {
		t.Fatalf("expected %s, got %s", expected2, got2)
	}
}

func TestRoutingVariable(t *testing.T) {
	var (
		expected = "variable"
		got      string
		mux      = New()
		w        = httptest.NewRecorder()
	)
	mux.Get("/:vartest", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = GetValue(r, "vartest")
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

func TestStaticFile(t *testing.T) {
	var data string
	mux := New()
	mux.Get("/file/", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		data = "DATA"
	}))

	r, _ := http.NewRequest("GET", "/file/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if data != "DATA" {
		t.Error("Data not serve")
	}
}

func TestStandAloneRoute(t *testing.T) {
	valid := false
	mux := http.NewServeMux()

	testRoute := NewRoute("/test", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		valid = true
	}))
	mux.Handle("/test", testRoute.Get())
	r, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !valid {
		t.Error("Route Handler not call")
	}
}

func TestRegexParam(t *testing.T) {
	valid := false
	mux := New()

	mux.Get("/Regex/#ttt^[a-z]$", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		valid = true
	}))

	r, _ := http.NewRequest("GET", "/Regex/test", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !valid {
		t.Error("Route Handler not call")
	}
}

func TestRegexParam2(t *testing.T) {
	valid := false
	mux := New()

	mux.Get("/Regex/#tttt^[a-z]$", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		valid = true
	}))

	r, _ := http.NewRequest("GET", "/Regex/1234", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if valid {
		t.Error("Regex param not valid !")
	}
}

func TestRegexParamMutli(t *testing.T) {
	valid := false
	mux := New()

	mux.Get("/Regex/#ttt^[a-z]$/#yyy^[0-9]$", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		valid = true
	}))

	r, _ := http.NewRequest("GET", "/Regex/first/2", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !valid {
		t.Error("Regex multi Params not valid !")
	}
}

func TestMultiParams(t *testing.T) {
	valid := false
	mux := New()

	mux.Get("/Regex/#num^[a-z]$/:test", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		valid = true
	}))

	r, _ := http.NewRequest("GET", "/Regex/first/second", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !valid {
		t.Error("Regex multi Params not valid !")
	}
}

func TestWC(t *testing.T) {
	valid := false
	mux := New()
	mux.GetFunc("/test/*", func(rw http.ResponseWriter, req *http.Request) {
		valid = true
	})

	req, _ := http.NewRequest("GET", "/test/random/route", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if !valid {
		t.Error("WC doesn't work !")
	}
}

func TestSlashRemoving(t *testing.T) {
	valid := false
	mux := New()
	mux.GetFunc("/test", func(rw http.ResponseWriter, req *http.Request) {
		valid = true
	})

	req, _ := http.NewRequest("GET", "/test/////", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if !valid {
		t.Error("Slash removing doesn't work !")
	}
}

func TestSubRouter(t *testing.T) {
	valid := false
	router := New()
	mux := New()

	mux.GetFunc("/test", func(rw http.ResponseWriter, req *http.Request) {
		valid = true
	})

	req, err := http.NewRequest("GET", "/index/test", nil)
	if err != nil {
		t.Error(err)
	}
	rw := httptest.NewRecorder()

	router.SubRoute("/index", mux)

	router.ServeHTTP(rw, req)

	if !valid {
		t.Error("sub route doesn't work !")
	}
}
