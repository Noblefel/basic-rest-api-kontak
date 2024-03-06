package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/models"
)

func TestNewUserHandlers(t *testing.T) {
	var db sql.DB
	user := NewUserHandlers(&db)

	typeString := reflect.TypeOf(user).String()

	if typeString != "*handlers.UserHandlers" {
		t.Error("NewUserHandlers() did not get the correct type, wanted *handlers.UserHandlers")
	}
}

func TestUser_All(t *testing.T) {
	var tests = []struct {
		name       string
		statusCode int
	}{
		{"success", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest("GET", "/users", nil)
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(h.user.All)
			handler.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("want %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestUser_Get(t *testing.T) {
	var tests = []struct {
		name       string
		userId     int
		statusCode int
	}{
		{"success", 1, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest("GET", fmt.Sprint("/users/", tt.userId), nil)
			ctx := context.WithValue(r.Context(), "user", models.User{Id: tt.userId})
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(h.user.Get)
			handler.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("want %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestUser_Update(t *testing.T) {
	var tests = []struct {
		name       string
		userId     int
		form       url.Values
		statusCode int
	}{
		{"success", 1, url.Values{"email": {"test@example.com"}, "password": {"password"}}, http.StatusOK},
		{"error parsing form", 1, nil, http.StatusBadRequest},
		{"error validation", 1, url.Values{"email": {"x"}}, http.StatusBadRequest},
		{"error updating user", -1, url.Values{"email": {"test@example.com"}, "password": {"password"}}, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r *http.Request
			if tt.form != nil {
				r, _ = http.NewRequest("PUT", fmt.Sprint("/users/", tt.userId), strings.NewReader(tt.form.Encode()))
			} else {
				r, _ = http.NewRequest("PUT", fmt.Sprint("/users/", tt.userId), nil)
			}

			ctx := context.WithValue(r.Context(), "user", models.User{Id: tt.userId})
			r = r.WithContext(ctx)
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(h.user.Update)
			handler.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("want %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestUser_Delete(t *testing.T) {
	var tests = []struct {
		name       string
		userId     int
		statusCode int
	}{
		{"success", 1, http.StatusOK},
		{"error deleting user", -1, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest("DELETE", fmt.Sprint("/users/", tt.userId), nil)
			ctx := context.WithValue(r.Context(), "user", models.User{Id: tt.userId})
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(h.user.Delete)
			handler.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("want %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}
