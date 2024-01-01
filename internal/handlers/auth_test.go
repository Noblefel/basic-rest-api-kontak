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

var authRegisterTests = []struct {
	name       string
	form       url.Values
	statusCode int
}{
	{
		name: "authRegister-ok",
		form: url.Values{
			"email":    {"test@example.com"},
			"password": {"password"},
		},
		statusCode: http.StatusCreated,
	},
	{
		name:       "authRegister-error-parsing-form",
		form:       nil,
		statusCode: http.StatusBadRequest,
	},
	{
		name: "authRegister-error-form-validation",
		form: url.Values{
			"email": {"x"},
		},
		statusCode: http.StatusBadRequest,
	},
	{
		name: "authRegister-error-duplicate-email",
		form: url.Values{
			"email":    {"alreadyexists@error.com"},
			"password": {"password"},
		},
		statusCode: http.StatusConflict,
	},
	{
		name: "authRegister-error-registering-user",
		form: url.Values{
			"email":    {"unexpected@error.com"},
			"password": {"password"},
		},
		statusCode: http.StatusInternalServerError,
	},
}

func TestAuth_Register(t *testing.T) {
	for _, tt := range authRegisterTests {
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
			t.Errorf("%s returned response code of %d, wanted %d", tt.name, w.Code, tt.statusCode)
		}
	}
}

var authLoginTests = []struct {
	name       string
	form       url.Values
	statusCode int
}{
	{
		name: "authLogin-ok",
		form: url.Values{
			"email":    {"test@example.com"},
			"password": {"password"},
		},
		statusCode: http.StatusOK,
	},
	{
		name:       "authLogin-error-parsing-form",
		form:       nil,
		statusCode: http.StatusBadRequest,
	},
	{
		name: "authLogin-error-form-validation",
		form: url.Values{
			"email": {"x"},
		},
		statusCode: http.StatusBadRequest,
	},
	{
		name: "authLogin-error-invalid-credentials",
		form: url.Values{
			"email":    {"test@example.com"},
			"password": {"wrongpassword"},
		},
		statusCode: http.StatusUnauthorized,
	},
	{
		name: "authLogin-error-authenticating",
		form: url.Values{
			"email":    {"test@example.com"},
			"password": {"unexpected error"},
		},
		statusCode: http.StatusInternalServerError,
	},
	{
		name: "authLogin-error-generating-jwt",
		form: url.Values{
			"email":    {"test@example.com"},
			"password": {"jwt error"},
		},
		statusCode: http.StatusInternalServerError,
	},
}

func TestAuth_Login(t *testing.T) {
	for _, tt := range authLoginTests {
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
			t.Errorf("%s returned response code of %d, wanted %d", tt.name, w.Code, tt.statusCode)
		}
	}
}
