package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestNewAuthHandlers(t *testing.T) {
	var db sql.DB
	auth := NewAuthHandlers(&db)

	typeString := reflect.TypeOf(auth).String()

	if typeString != "*handlers.AuthHandlers" {
		t.Error("NewAuthHandlers() did not get the correct type, wanted *handlers.AuthHandlers")
	}
}

func TestAuth_Register(t *testing.T) {
	var tests = []struct {
		name       string
		form       url.Values
		statusCode int
	}{
		{"success", url.Values{"email": {"test@example.com"}, "password": {"password"}}, http.StatusCreated},
		{"error parsing form", nil, http.StatusBadRequest},
		{"error validation", url.Values{"email": {"x"}}, http.StatusBadRequest},
		{"duplicate email", url.Values{"email": {"alreadyexists@error.com"}, "password": {"password"}}, http.StatusConflict},
		{"error registering ", url.Values{"email": {"unexpected@error.com"}, "password": {"password"}}, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r *http.Request
			if tt.form != nil {
				r, _ = http.NewRequest("POST", "/auth/register", strings.NewReader(tt.form.Encode()))
			} else {
				r, _ = http.NewRequest("POST", "/auth/register", nil)
			}

			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(h.auth.Register)
			handler.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("want %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestAuth_Login(t *testing.T) {
	var tests = []struct {
		name       string
		form       url.Values
		statusCode int
	}{
		{"success", url.Values{"email": {"test@example.com"}, "password": {"password"}}, http.StatusOK},
		{"error parsing form", nil, http.StatusBadRequest},
		{"error validation", url.Values{"email": {"x"}}, http.StatusBadRequest},
		{"invalid credentials", url.Values{"email": {"test@example.com"}, "password": {"wrongpassword"}}, http.StatusUnauthorized},
		{"error authenticating", url.Values{"email": {"test@example.com"}, "password": {"unexpected error"}}, http.StatusInternalServerError},
		{"error generating jwt", url.Values{"email": {"test@example.com"}, "password": {"jwt error"}}, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r *http.Request
			if tt.form != nil {
				r, _ = http.NewRequest("POST", "/auth/login", strings.NewReader(tt.form.Encode()))
			} else {
				r, _ = http.NewRequest("POST", "/auth/login", nil)
			}

			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(h.auth.Login)
			handler.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("want %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}
