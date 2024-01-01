package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestForm_New(t *testing.T) {
	form := New(url.Values{})

	typeString := reflect.TypeOf(form).String()
	if typeString != "*forms.Form" {
		t.Error("New() did not get the correct type, wanted *forms.Form")
	}
}

func TestFormErrors_Add(t *testing.T) {
	form := New(url.Values{})
	form.Errors.Add("field", "message")
	form.Errors.Add("field2", "message")

	if len(form.Errors) == 0 {
		t.Error("Errors Add() did not add a new error")
	}
}

func TestForm_Valid(t *testing.T) {
	form := New(url.Values{})

	if !form.Valid() {
		t.Error("Form should have been valid")
	}

	form.Errors.Add("field", "some error")

	if form.Valid() {
		t.Error("Form should have been invalid")
	}
}

func TestForm_Required(t *testing.T) {
	form := New(url.Values{})
	form.Required("name", "email")

	if form.Valid() {
		t.Error("Form is valid despite the required fields")
	}

	form = New(url.Values{
		"name": {"test"},
	})
	form.Required("name")

	if !form.Valid() {
		t.Error("Form is valid when required fields are fulfilled")
	}
}

func TestForm_StringMinLength(t *testing.T) {
	form := New(url.Values{
		"password": {""},
	})
	form.StringMinLength("password", 8)

	if !form.Valid() {
		t.Error("Form string min length should return valid when field is empty")
	}

	form = New(url.Values{
		"password": {"x"},
	})
	form.StringMinLength("password", 8)

	if form.Valid() {
		t.Error("Form string min length returns valid when the value is less")
	}

	form = New(url.Values{
		"password": {"12345678910"},
	})
	form.StringMinLength("password", 8)

	if !form.Valid() {
		t.Error("Form string min length returns invalid when the length is longer")
	}
}

func TestForm_Email(t *testing.T) {
	form := New(url.Values{
		"email": {""},
	})
	form.Email("email")

	if !form.Valid() {
		t.Error("Form email should return valid when field is empty")
	}

	form = New(url.Values{
		"email": {"x"},
	})
	form.Email("email")

	if form.Valid() {
		t.Error("Form email should return invalid for incorrect email")
	}

	form = New(url.Values{
		"email": {"test@example.com"},
	})
	form.Email("email")

	if !form.Valid() {
		t.Error("Form email should return valid for correct email")
	}
}

func TestForm_ValidOrErr(t *testing.T) {
	form := New(url.Values{})
	form.Required("field")

	r, _ := http.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()

	if form.ValidOrErr(w, r) {
		t.Error("Form ValidOrErr should return false for invalid form")
	}

	form = New(url.Values{})
	if !form.ValidOrErr(w, r) {
		t.Error("Form ValidOrErr should return true for valid form")
	}
}
