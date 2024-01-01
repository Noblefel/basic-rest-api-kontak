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

var contactAllTests = []struct {
	name       string
	statusCode int
}{
	{
		name:       "contactAll-ok",
		statusCode: http.StatusOK,
	},
}

func TestContact_All(t *testing.T) {
	for _, tt := range contactAllTests {
		r, _ := http.NewRequest("GET", "/contacts", nil)
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(h.contact.All)
		handler.ServeHTTP(w, r)

		if w.Code != tt.statusCode {
			t.Errorf("%s returned response code of %d, wanted %d", tt.name, w.Code, tt.statusCode)
		}
	}
}

var contactGetTests = []struct {
	name       string
	contact    models.Contact
	statusCode int
}{
	{
		name: "contactAll-ok",
		contact: models.Contact{
			Id: 1,
		},
		statusCode: http.StatusOK,
	},
}

func TestContact_Get(t *testing.T) {
	for _, tt := range contactGetTests {
		r, _ := http.NewRequest("GET", fmt.Sprint("/contacts/", tt.contact.Id), nil)
		ctx := context.WithValue(r.Context(), "contact", tt.contact)
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(h.contact.Get)
		handler.ServeHTTP(w, r)

		if w.Code != tt.statusCode {
			t.Errorf("%s returned response code of %d, wanted %d", tt.name, w.Code, tt.statusCode)
		}
	}
}

var contactGetByUserTests = []struct {
	name       string
	user       models.User
	statusCode int
}{
	{
		name: "contactGetByUser-ok",
		user: models.User{
			Id: 1,
		},
		statusCode: http.StatusOK,
	},
	{
		name: "contactGetByUser-error-getting-user-contacts",
		user: models.User{
			Id: -1,
		},
		statusCode: http.StatusInternalServerError,
	},
}

func TestContact_GetByUser(t *testing.T) {
	for _, tt := range contactGetByUserTests {
		r, _ := http.NewRequest("GET", fmt.Sprint("/users/", tt.user.Id, "/contacts"), nil)
		ctx := context.WithValue(r.Context(), "user", tt.user)
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(h.contact.GetByUser)
		handler.ServeHTTP(w, r)

		if w.Code != tt.statusCode {
			t.Errorf("%s returned response code of %d, wanted %d", tt.name, w.Code, tt.statusCode)
		}
	}
}

var contactCreateTests = []struct {
	name       string
	userId     int
	form       url.Values
	statusCode int
}{
	{
		name:   "contactCreate-ok",
		userId: 1,
		form: url.Values{
			"nama": {"John Doe"},
		},
		statusCode: http.StatusCreated,
	},
	{
		name:       "contactCreate-error-parsing-form",
		userId:     1,
		form:       nil,
		statusCode: http.StatusBadRequest,
	},
	{
		name:   "contactCreate-error-form-validation",
		userId: 1,
		form: url.Values{
			"nama": {""},
		},
		statusCode: http.StatusBadRequest,
	},
	{
		name:   "contactCreate-error-creating-contact",
		userId: -1,
		form: url.Values{
			"nama": {"John Doe"},
		},
		statusCode: http.StatusInternalServerError,
	},
}

func TestContact_Create(t *testing.T) {
	for _, tt := range contactCreateTests {
		var r *http.Request
		if tt.form != nil {
			r, _ = http.NewRequest("POST", "/contacts/create", strings.NewReader(tt.form.Encode()))
		} else {
			r, _ = http.NewRequest("POST", "/contacts/create", nil)
		}

		ctx := context.WithValue(r.Context(), "user_id", tt.userId)
		r = r.WithContext(ctx)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(h.contact.Create)
		handler.ServeHTTP(w, r)

		if w.Code != tt.statusCode {
			t.Errorf("%s returned response code of %d, wanted %d", tt.name, w.Code, tt.statusCode)
		}
	}
}

var contactUpdateTests = []struct {
	name       string
	contact    models.Contact
	form       url.Values
	statusCode int
}{
	{
		name:    "contactUpdate-ok",
		contact: models.Contact{Id: 1},
		form: url.Values{
			"nama": {"John Doe"},
		},
		statusCode: http.StatusOK,
	},
	{
		name:       "contactUpdate-error-parsing-form",
		contact:    models.Contact{Id: 1},
		form:       nil,
		statusCode: http.StatusBadRequest,
	},
	{
		name:    "contactUpdate-error-form-validation",
		contact: models.Contact{Id: 1},
		form: url.Values{
			"nama": {""},
		},
		statusCode: http.StatusBadRequest,
	},
	{
		name:    "contactUpdate-error-updating-contact",
		contact: models.Contact{Id: -1},
		form: url.Values{
			"nama": {"John Doe"},
		},
		statusCode: http.StatusInternalServerError,
	},
}

func TestContact_Update(t *testing.T) {
	for _, tt := range contactUpdateTests {
		var r *http.Request
		if tt.form != nil {
			r, _ = http.NewRequest("POST", fmt.Sprint("/contacts/", tt.contact.Id, "/update"), strings.NewReader(tt.form.Encode()))
		} else {
			r, _ = http.NewRequest("POST", fmt.Sprint("/contacts/", tt.contact.Id, "/update"), nil)
		}

		ctx := context.WithValue(r.Context(), "contact", tt.contact)
		r = r.WithContext(ctx)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(h.contact.Update)
		handler.ServeHTTP(w, r)

		if w.Code != tt.statusCode {
			t.Errorf("%s returned response code of %d, wanted %d", tt.name, w.Code, tt.statusCode)
		}
	}
}

var contactDeleteTests = []struct {
	name       string
	contact    models.Contact
	statusCode int
}{
	{
		name:       "contactDelete-ok",
		contact:    models.Contact{Id: 1},
		statusCode: http.StatusOK,
	},
	{
		name:       "contactDelete-error-deleting-contact",
		contact:    models.Contact{Id: -1},
		statusCode: http.StatusInternalServerError,
	},
}

func TestContact_Delete(t *testing.T) {
	for _, tt := range contactDeleteTests {
		r, _ := http.NewRequest("POST", fmt.Sprint("/contacts/", tt.contact.Id, "/delete"), nil)
		ctx := context.WithValue(r.Context(), "contact", tt.contact)
		r = r.WithContext(ctx)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(h.contact.Delete)
		handler.ServeHTTP(w, r)

		if w.Code != tt.statusCode {
			t.Errorf("%s returned response code of %d, wanted %d", tt.name, w.Code, tt.statusCode)
		}
	}
}
