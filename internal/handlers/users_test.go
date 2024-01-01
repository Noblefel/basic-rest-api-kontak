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

var userAllTests = []struct {
	name       string
	statusCode int
}{
	{
		name:       "userAll-ok",
		statusCode: http.StatusOK,
	},
}

func TestUser_All(t *testing.T) {
	for _, tt := range userAllTests {
		r, _ := http.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(h.user.All)
		handler.ServeHTTP(w, r)

		if w.Code != tt.statusCode {
			t.Errorf("%s returned response code of %d, wanted %d", tt.name, w.Code, tt.statusCode)
		}
	}
}

var userGetTests = []struct {
	name       string
	user       models.User
	statusCode int
}{
	{
		name: "userAll-ok",
		user: models.User{
			Id: 1,
		},
		statusCode: http.StatusOK,
	},
}

func TestUser_Get(t *testing.T) {
	for _, tt := range userGetTests {
		r, _ := http.NewRequest("GET", fmt.Sprint("/users/", tt.user.Id), nil)
		ctx := context.WithValue(r.Context(), "user", tt.user)
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(h.user.Get)
		handler.ServeHTTP(w, r)

		if w.Code != tt.statusCode {
			t.Errorf("%s returned response code of %d, wanted %d", tt.name, w.Code, tt.statusCode)
		}
	}
}

var userUpdateTests = []struct {
	name       string
	user       models.User
	form       url.Values
	statusCode int
}{
	{
		name: "userUpdate-ok",
		user: models.User{Id: 1},
		form: url.Values{
			"email":    {"test@example.com"},
			"password": {"password"},
		},
		statusCode: http.StatusOK,
	},
	{
		name:       "userUpdate-error-parsing-form",
		user:       models.User{},
		form:       nil,
		statusCode: http.StatusBadRequest,
	},
	{
		name: "userUpdate-error-form-validation",
		user: models.User{},
		form: url.Values{
			"email": {"x"},
		},
		statusCode: http.StatusBadRequest,
	},
	{
		name: "userUpdate-error-updating-user",
		user: models.User{Id: -1},
		form: url.Values{
			"email":    {"test@example.com"},
			"password": {"password"},
		},
		statusCode: http.StatusInternalServerError,
	},
}

func TestUser_Update(t *testing.T) {
	for _, tt := range userUpdateTests {
		var r *http.Request
		if tt.form != nil {
			r, _ = http.NewRequest("POST", fmt.Sprint("/users/", tt.user.Id), strings.NewReader(tt.form.Encode()))
		} else {
			r, _ = http.NewRequest("POST", fmt.Sprint("/users/", tt.user.Id), nil)
		}

		ctx := context.WithValue(r.Context(), "user", tt.user)
		r = r.WithContext(ctx)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(h.user.Update)
		handler.ServeHTTP(w, r)

		if w.Code != tt.statusCode {
			t.Errorf("%s returned response code of %d, wanted %d", tt.name, w.Code, tt.statusCode)
		}
	}
}

var userDeleteTests = []struct {
	name       string
	user       models.User
	statusCode int
}{
	{
		name:       "userDelete-ok",
		user:       models.User{Id: 1},
		statusCode: http.StatusOK,
	},
	{
		name:       "userDelete-error-deleting-user",
		user:       models.User{Id: -1},
		statusCode: http.StatusInternalServerError,
	},
}

func TestUser_Delete(t *testing.T) {
	for _, tt := range userDeleteTests {
		r, _ := http.NewRequest("POST", fmt.Sprint("/users/", tt.user.Id), nil)
		ctx := context.WithValue(r.Context(), "user", tt.user)
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(h.user.Delete)
		handler.ServeHTTP(w, r)

		if w.Code != tt.statusCode {
			t.Errorf("%s returned response code of %d, wanted %d", tt.name, w.Code, tt.statusCode)
		}
	}
}
