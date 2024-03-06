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

func TestNewContactHandlers(t *testing.T) {
	var db sql.DB
	contact := NewContactHandlers(&db)

	typeString := reflect.TypeOf(contact).String()

	if typeString != "*handlers.ContactHandlers" {
		t.Error("NewContactHandlers() did not get the correct type, wanted *handlers.ContactHandlers")
	}
}

func TestContact_All(t *testing.T) {
	var tests = []struct {
		name       string
		statusCode int
	}{
		{"success", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest("GET", "/contacts", nil)
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(h.contact.All)
			handler.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("want %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestContact_Get(t *testing.T) {
	var tests = []struct {
		name       string
		contactId  int
		statusCode int
	}{
		{"success", 1, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest("GET", fmt.Sprint("/contacts/", tt.contactId), nil)
			ctx := context.WithValue(r.Context(), "contact", models.Contact{Id: tt.contactId})
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(h.contact.Get)
			handler.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("want %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestContact_GetByUser(t *testing.T) {
	var tests = []struct {
		name       string
		user       models.User
		statusCode int
	}{
		{"success", models.User{Id: 1}, http.StatusOK},
		{"error getting contacts", models.User{Id: -1}, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest("GET", fmt.Sprint("/users/", tt.user.Id, "/contacts"), nil)
			ctx := context.WithValue(r.Context(), "user", tt.user)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(h.contact.GetByUser)
			handler.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("want %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestContact_Create(t *testing.T) {
	var tests = []struct {
		name       string
		userId     int
		form       url.Values
		statusCode int
	}{
		{"success", 1, url.Values{"nama": {"John Doe"}}, http.StatusCreated},
		{"error parsing form", 1, nil, http.StatusBadRequest},
		{"error validation", 1, url.Values{"nama": {""}}, http.StatusBadRequest},
		{"error creating contact", -1, url.Values{"nama": {"John Doe"}}, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r *http.Request
			if tt.form != nil {
				r, _ = http.NewRequest("POST", "/contacts", strings.NewReader(tt.form.Encode()))
			} else {
				r, _ = http.NewRequest("POST", "/contacts", nil)
			}

			ctx := context.WithValue(r.Context(), "user_id", tt.userId)
			r = r.WithContext(ctx)
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(h.contact.Create)
			handler.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("want %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestContact_Update(t *testing.T) {
	var tests = []struct {
		name       string
		contactId  int
		form       url.Values
		statusCode int
	}{
		{"success", 1, url.Values{"nama": {"John Doe"}}, http.StatusOK},
		{"error parsing form", 1, nil, http.StatusBadRequest},
		{"error validation", 1, url.Values{"nama": {""}}, http.StatusBadRequest},
		{"error updating contact", -1, url.Values{"nama": {"John Doe"}}, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r *http.Request
			if tt.form != nil {
				r, _ = http.NewRequest("PUT", fmt.Sprint("/contacts/", tt.contactId), strings.NewReader(tt.form.Encode()))
			} else {
				r, _ = http.NewRequest("PUT", fmt.Sprint("/contacts/", tt.contactId), nil)
			}

			ctx := context.WithValue(r.Context(), "contact", models.Contact{Id: tt.contactId})
			r = r.WithContext(ctx)
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(h.contact.Update)
			handler.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("want %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestContact_Delete(t *testing.T) {
	var tests = []struct {
		name       string
		contactId  int
		statusCode int
	}{
		{"success", 1, http.StatusOK},
		{"error deleting contact", -1, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest("DELETE", fmt.Sprint("/contacts/", tt.contactId), nil)
			ctx := context.WithValue(r.Context(), "contact", models.Contact{Id: tt.contactId})
			r = r.WithContext(ctx)
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(h.contact.Delete)
			handler.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("want %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}
