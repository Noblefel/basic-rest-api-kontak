package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type testHandlers struct {
	auth    *AuthHandlers
	user    *UserHandlers
	contact *ContactHandlers
}

func newTestHandlers() *testHandlers {
	return &testHandlers{
		auth:    NewTestAuthHandlers(),
		user:    NewTestUserHandlers(),
		contact: NewTestContactHandlers(),
	}
}

var h = newTestHandlers()

func TestBaseHandlers(t *testing.T) {
	var tests = []struct {
		name       string
		url        string
		method     string
		statusCode int
	}{
		{"not-found", "/xmo02v3o2cm3ro", "GET", http.StatusNotFound},
		{"method-not-allowed", "/users", "POST", http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()

			var handler http.HandlerFunc
			if tt.name == "not-found" {
				handler = NotFound
			} else {
				handler = MethodNotAllowed
			}

			handler.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("want %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}
