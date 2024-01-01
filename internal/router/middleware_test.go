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

var middlewareAuthTests = []struct {
	name           string
	authorization  string
	expectedUserId int
	statusCode     int
}{
	{
		name:           "middlewareAuth-ok",
		authorization:  sampleToken,
		expectedUserId: 1,
		statusCode:     http.StatusOK,
	},
	{
		name:           "middlewareAuth-ok-2",
		authorization:  sampleToken2,
		expectedUserId: 2,
		statusCode:     http.StatusOK,
	},
	{
		name:          "middlewareAuth-empty-authorization-header",
		authorization: "",
		statusCode:    http.StatusUnauthorized,
	},
	{
		name:          "middlewareAuth-invalid-authorization-header",
		authorization: "abcdefghijklmn",
		statusCode:    http.StatusUnauthorized,
	},
}

func TestMiddleware_Auth(t *testing.T) {
	for _, tt := range middlewareAuthTests {
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
			t.Errorf("%s returned response code of %d, wanted %d", tt.name, w.Code, tt.statusCode)
		}
	}
}

var middlewareUserGuardTests = []struct {
	name        string
	userId      int
	userLevel   int
	userIdRoute string
	statusCode  int
}{
	{
		name:        "middlewareUserGuard-ok",
		userId:      1,
		userLevel:   1,
		userIdRoute: "2",
		statusCode:  http.StatusOK,
	},
	{
		name:        "middlewareUserGuard-error-invalid-id",
		userId:      1,
		userLevel:   1,
		userIdRoute: "invalid",
		statusCode:  http.StatusBadRequest,
	},
	{
		name:        "middlewareUserGuard-error-unauthorized",
		userId:      1,
		userLevel:   0,
		userIdRoute: "2",
		statusCode:  http.StatusUnauthorized,
	},
	{
		name:        "middlewareUserGuard-error-user-not-found",
		userId:      1,
		userLevel:   1,
		userIdRoute: "999999999",
		statusCode:  http.StatusNotFound,
	},
	{
		name:        "middlewareUserGuard-error-getting-user",
		userId:      1,
		userLevel:   1,
		userIdRoute: "-1",
		statusCode:  http.StatusInternalServerError,
	},
}

func TestMiddleware_UserGuard(t *testing.T) {
	for _, tt := range middlewareUserGuardTests {
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
			t.Errorf("%s returned response code of %d, wanted %d", tt.name, w.Code, tt.statusCode)
		}
	}
}

var middlewareContactGuardTests = []struct {
	name           string
	userId         int
	userLevel      int
	contactIdRoute string
	statusCode     int
}{
	{
		name:           "middlewareContactGuard-ok",
		userId:         1,
		userLevel:      1,
		contactIdRoute: "1",
		statusCode:     http.StatusOK,
	},
	{
		name:           "middlewareContactGuard-error-invalid-id",
		userId:         1,
		userLevel:      1,
		contactIdRoute: "invalid",
		statusCode:     http.StatusBadRequest,
	},
	{
		name:           "middlewareContactGuard-error-contact-not-found",
		userId:         1,
		userLevel:      1,
		contactIdRoute: "99999999",
		statusCode:     http.StatusNotFound,
	},
	{
		name:           "middlewareContactGuard-error-getting-contact",
		userId:         1,
		userLevel:      1,
		contactIdRoute: "-1",
		statusCode:     http.StatusInternalServerError,
	},
	{
		name:           "middlewareContactGuard-error-unauthorized",
		userId:         2,
		userLevel:      0,
		contactIdRoute: "1",
		statusCode:     http.StatusUnauthorized,
	},
}

func TestMiddleware_ContactGuard(t *testing.T) {
	for _, tt := range middlewareContactGuardTests {
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
			t.Errorf("%s returned response code of %d, wanted %d", tt.name, w.Code, tt.statusCode)
		}
	}
}

var middlewareAdminOnlyTests = []struct {
	name       string
	userLevel  int
	statusCode int
}{
	{
		name:       "middlewareAdminOnly-ok",
		userLevel:  1,
		statusCode: http.StatusOK,
	},
	{
		name:       "middlewareAdminOnly-error-unauthorized",
		userLevel:  0,
		statusCode: http.StatusUnauthorized,
	},
}

func TestMiddleware_AdminOnly(t *testing.T) {
	for _, tt := range middlewareAdminOnlyTests {
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

		r := httptest.NewRequest("GET", "/", nil)
		ctx := context.WithValue(r.Context(), "level", tt.userLevel)
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()

		h := m.AdminOnly(next)
		h.ServeHTTP(w, r)

		if w.Code != tt.statusCode {
			t.Errorf("%s returned response code of %d, wanted %d", tt.name, w.Code, tt.statusCode)
		}
	}
}
