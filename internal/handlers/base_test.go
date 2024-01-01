package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var baseTests = []struct {
	name       string
	url        string
	method     string
	statusCode int
}{
	{
		name:       "not-found",
		url:        "/xmo02v3o2cm3ro",
		method:     "GET",
		statusCode: http.StatusNotFound,
	},
	{
		name:       "method-not-allowed",
		url:        "/users",
		method:     "POST",
		statusCode: http.StatusMethodNotAllowed,
	},
}

func TestBaseHandlers(t *testing.T) {
	for _, tt := range baseTests {
		r, _ := http.NewRequest(tt.method, tt.url, nil)
		w := httptest.NewRecorder()

		var handler http.HandlerFunc
		if tt.name == "not-found" {
			handler = NotFound
		} else {
			handler = MethodNotAllowed
		}

		handler.ServeHTTP(w, r)

		if w.Code != tt.statusCode {
			t.Errorf("%s returned response code of %d, wanted %d", tt.name, w.Code, tt.statusCode)
		}
	}
}
