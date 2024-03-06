package router

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	u "github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/utils"
	"github.com/go-chi/chi/v5"
)

type Params map[string]string

func getCtxWithParam(r *http.Request, params Params) context.Context {
	ctx := r.Context()
	chiCtx := chi.NewRouteContext()
	for k, v := range params {
		chiCtx.URLParams.Add(k, v)
	}
	ctx = context.WithValue(ctx, chi.RouteCtxKey, chiCtx)
	return ctx
}

func TestNewMiddleware(t *testing.T) {
	var db *sql.DB
	middleware := NewMiddleware(db)

	typeString := reflect.TypeOf(middleware).String()
	if typeString != "*router.Middleware" {
		t.Error("NewMiddleware() did not get the correct type, wanted *router.Middleware")
	}
}

var m = NewTestMiddleware()

var sampleToken, _ = u.GenerateJWT(1, 1)
var sampleToken2, _ = u.GenerateJWT(2, 1)

func TestMiddleware_Auth(t *testing.T) {
	var tests = []struct {
		name           string
		authorization  string
		expectedUserId int
		statusCode     int
	}{
		{"success", sampleToken, 1, http.StatusOK},
		{"success 2", sampleToken2, 2, http.StatusOK},
		{"empty authorization header", "", 0, http.StatusUnauthorized},
		{"invalid authorization header", "abcdefghijklmn", 0, http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userId := r.Context().Value("user_id")
				if userId == nil {
					t.Error("User id not in context")
					return
				}

				if userId.(int) != tt.expectedUserId {
					t.Error("Expected user id does not match")
					return
				}
			})

			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Set("Authorization", tt.authorization)
			w := httptest.NewRecorder()
			h := m.Auth(next)
			h.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("want %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestMiddleware_UserGuard(t *testing.T) {
	var tests = []struct {
		name        string
		userId      int
		userLevel   int
		userIdRoute string
		statusCode  int
	}{
		{"success", 1, 1, "2", http.StatusOK},
		{"invalid id", 1, 1, "invalid", http.StatusBadRequest},
		{"unauthorized", 1, 0, "2", http.StatusUnauthorized},
		{"user not found", 1, 1, "999999999", http.StatusNotFound},
		{"error getting user", 1, 1, "-1", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			r := httptest.NewRequest("GET", "/{user_id}", nil)
			ctx := getCtxWithParam(r, Params{"user_id": tt.userIdRoute})
			ctx = context.WithValue(ctx, "user_id", tt.userId)
			ctx = context.WithValue(ctx, "level", tt.userLevel)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()
			h := m.UserGuard(next)
			h.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("want %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestMiddleware_ContactGuard(t *testing.T) {
	var tests = []struct {
		name           string
		userId         int
		userLevel      int
		contactIdRoute string
		statusCode     int
	}{
		{"success", 1, 1, "1", http.StatusOK},
		{"invalid id", 1, 1, "invalid", http.StatusBadRequest},
		{"contact not found", 1, 1, "99999999", http.StatusNotFound},
		{"error getting contact", 1, 1, "-1", http.StatusInternalServerError},
		{"unauthorized", 2, 0, "1", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			r := httptest.NewRequest("GET", "/{contact_id}", nil)
			ctx := getCtxWithParam(r, Params{"contact_id": tt.contactIdRoute})
			ctx = context.WithValue(ctx, "user_id", tt.userId)
			ctx = context.WithValue(ctx, "level", tt.userLevel)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()
			h := m.ContactGuard(next)
			h.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("want %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestMiddleware_AdminOnly(t *testing.T) {
	var tests = []struct {
		name       string
		userLevel  int
		statusCode int
	}{
		{"success", 1, http.StatusOK},
		{"unauthorized", 0, http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			r := httptest.NewRequest("GET", "/", nil)
			ctx := context.WithValue(r.Context(), "level", tt.userLevel)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()
			h := m.AdminOnly(next)
			h.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("want %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}
