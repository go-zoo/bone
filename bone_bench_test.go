package bone

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daryl/zeus"
	"github.com/gorilla/mux"
	"github.com/gorilla/pat"
	"github.com/julienschmidt/httprouter"
)

// Test the ns/op
func BenchmarkBoneMux(b *testing.B) {
	request, _ := http.NewRequest("GET", "/sd/df", nil)
	response := httptest.NewRecorder()
	muxx := New()

	muxx.Get("/", http.HandlerFunc(Bench))
	muxx.Get("/aas", http.HandlerFunc(Bench))
	muxx.Get("/aasr", http.HandlerFunc(Bench))
	muxx.Get("/sd", http.HandlerFunc(Bench))
	muxx.Get("/sd/:p", http.HandlerFunc(Bench))
	muxx.Get("/a", http.HandlerFunc(Bench))

	for n := 0; n < b.N; n++ {
		muxx.ServeHTTP(response, request)
	}
}

// Test daryl/zeus ns/op
func BenchmarkZeusMux(b *testing.B) {
	request, _ := http.NewRequest("GET", "/sd/df", nil)
	response := httptest.NewRecorder()
	muxx := zeus.New()

	muxx.GET("/", Bench)
	muxx.POST("/a", Bench)
	muxx.GET("/aas", Bench)
	muxx.GET("/sd/:p", Bench)

	for n := 0; n < b.N; n++ {
		muxx.ServeHTTP(response, request)
	}
}

// Test httprouter ns/op
func BenchmarkHttpRouterMux(b *testing.B) {
	request, _ := http.NewRequest("GET", "/sd/df", nil)
	response := httptest.NewRecorder()
	muxx := httprouter.New()
	muxx.Handler("GET", "/", http.HandlerFunc(Bench))
	muxx.Handler("POST", "/a", http.HandlerFunc(Bench))
	muxx.Handler("GET", "/aas", http.HandlerFunc(Bench))
	muxx.Handler("GET", "/sd/:p", http.HandlerFunc(Bench))

	for n := 0; n < b.N; n++ {
		muxx.ServeHTTP(response, request)
	}
}

// Test net/http ns/op
func BenchmarkNetHttpMux(b *testing.B) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	muxx := http.NewServeMux()
	muxx.HandleFunc("/", Bench)

	for n := 0; n < b.N; n++ {
		muxx.ServeHTTP(response, request)
	}
}

// Test gorilla/mux ns/op
func BenchmarkGorillaMux(b *testing.B) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	muxx := mux.NewRouter()
	muxx.Handle("/", http.HandlerFunc(Bench))

	for n := 0; n < b.N; n++ {
		muxx.ServeHTTP(response, request)
	}
}

// Test gorilla/pat ns/op
func BenchmarkGorillaPatMux(b *testing.B) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	muxx := pat.New()
	muxx.Get("/", Bench)

	for n := 0; n < b.N; n++ {
		muxx.ServeHTTP(response, request)
	}
}

func Bench(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("b"))
}

/*			### Result ###

BenchmarkBoneMux        10000000               118 ns/op
BenchmarkZeusMux          100000             54813 ns/op
BenchmarkHttpRouterMux  10000000               143 ns/op
BenchmarkNetHttpMux      3000000               548 ns/op
BenchmarkGorillaMux       300000              3333 ns/op
BenchmarkGorillaPatMux   1000000              1889 ns/op

ok      github.com/squiidz/bone 14.084s

*/
