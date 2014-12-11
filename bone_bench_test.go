package bone

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test the op/s
func BenchmarkMux(b *testing.B) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	muxx := New()

	for n := 0; n < b.N; n++ {
		muxx.ServeHTTP(response, request)
	}
}
