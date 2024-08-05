package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Noblefel/baic-rest-api-kontak/internal/models"
	"github.com/Noblefel/baic-rest-api-kontak/internal/storage"
	"github.com/Noblefel/baic-rest-api-kontak/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

var handlers *Handlers

func init() {
	store := storage.New()
	handlers = New(store)
}

func TestRegister(t *testing.T) {
	t.Cleanup(handlers.storage.Reset)

	t.Run("success", func(t *testing.T) {
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(models.Auth{
			Email:    "abc@example.com",
			Password: "12345678",
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/", &buf)
		req.Header.Add("Content-Type", "application/json")
		handlers.Register(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("want 201 created, got %d", rec.Code)
		}
	})

	t.Run("duplicate email", func(t *testing.T) {
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(models.Auth{
			Email:    "abc@example.com",
			Password: "12345678",
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/", &buf)
		req.Header.Add("Content-Type", "application/json")
		handlers.Register(rec, req)

		if rec.Code != http.StatusConflict {
			t.Errorf("want 409 conflict, got %d", rec.Code)
		}
	})
}

func TestLogin(t *testing.T) {
	t.Cleanup(handlers.storage.Reset)

	hash, _ := bcrypt.GenerateFromPassword([]byte("12345678"), 10)
	handlers.storage.Register("abc@example.com", string(hash))

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(models.Auth{
		Email:    "abc@example.com",
		Password: "12345678",
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", &buf)
	req.Header.Add("Content-Type", "application/json")
	handlers.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200 ok, got %d", rec.Code)
	}

	var resp utils.Response
	json.NewDecoder(rec.Body).Decode(&resp)

	if resp.Data == "" {
		t.Error("no token")
	}

	json.NewEncoder(&buf).Encode(models.Auth{
		Email:    "abc@example.com",
		Password: "asdonaosidnoasnd",
	})

	rec = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/", &buf)
	req.Header.Add("Content-Type", "application/json")
	handlers.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("want 401 unauthorized, got %d", rec.Code)
	}
}

func TestGetAccount(t *testing.T) {
	t.Cleanup(handlers.storage.Reset)
	handlers.storage.Register("abc", "")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), "user_id", 1)
	req = req.WithContext(ctx)
	handlers.GetAccount(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("want 200 ok, got %d", rec.Code)
	}

	rec = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/", nil)
	ctx = context.WithValue(req.Context(), "user_id", 2)
	req = req.WithContext(ctx)
	handlers.GetAccount(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("want 404 not found, got %d", rec.Code)
	}
}

func TestUpdateAccount(t *testing.T) {
	t.Cleanup(handlers.storage.Reset)
	handlers.storage.Register("abc", "")

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(models.Auth{
		Email: "new@example.com",
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/", &buf)
	ctx := context.WithValue(req.Context(), "user_id", 1)
	req = req.WithContext(ctx)
	handlers.UpdateAccount(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("want 200 ok, got %d", rec.Code)
	}
}

func TestDeleteAccount(t *testing.T) {
	t.Cleanup(handlers.storage.Reset)
	handlers.storage.Register("abc", "")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/", nil)
	ctx := context.WithValue(req.Context(), "user_id", 1)
	req = req.WithContext(ctx)
	handlers.DeleteAccount(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("want 200 ok, got %d", rec.Code)
	}
}

func TestGetAccountContacts(t *testing.T) {
	t.Cleanup(handlers.storage.Reset)
	handlers.storage.CreateContact(models.Contact{UserId: 1})
	handlers.storage.CreateContact(models.Contact{UserId: 1})
	handlers.storage.CreateContact(models.Contact{UserId: 2})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), "user_id", 1)
	req = req.WithContext(ctx)
	handlers.GetAccountContacts(rec, req)

	var body struct {
		Data []models.Contact `json:"data"`
	}

	err := json.NewDecoder(rec.Body).Decode(&body)
	if err != nil {
		t.Fatal(err)
	}

	if len(body.Data) != 2 {
		t.Errorf("expecting 2 contacts, got %d", len(body.Data))
	}
}

func TestGetContact(t *testing.T) {
	t.Cleanup(handlers.storage.Reset)
	handlers.storage.CreateContact(models.Contact{UserId: 1})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.SetPathValue("id", "1")
	ctx := context.WithValue(req.Context(), "user_id", 1)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	handlers.GetContact(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("want 200 ok, got %d", rec.Code)
	}

	rec = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/", nil)
	ctx = context.WithValue(req.Context(), "user_id", 2)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	handlers.GetContact(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("want 401 unauthorized, got %d", rec.Code)
	}
}

func TestCreateContact(t *testing.T) {
	t.Cleanup(handlers.storage.Reset)
	handlers.storage.CreateContact(models.Contact{})

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(models.Contact{
		Nama:  "abc",
		Email: "abc@example.com",
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", &buf)
	ctx := context.WithValue(req.Context(), "user_id", 1)
	req = req.WithContext(ctx)
	handlers.CreateContact(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("want 201 created, got %d", rec.Code)
	}

	var body struct {
		NewId int `json:"data"`
	}

	json.NewDecoder(rec.Body).Decode(&body)

	if body.NewId != 2 {
		t.Errorf("should return new id of 2, got %d", body.NewId)
	}
}

func TestUpdateContact(t *testing.T) {
	t.Cleanup(handlers.storage.Reset)
	handlers.storage.CreateContact(models.Contact{UserId: 1})

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(models.Contact{
		Nama:  "abc",
		Email: "abc@example.com",
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/", &buf)
	ctx := context.WithValue(req.Context(), "user_id", 1)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	handlers.UpdateContact(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("want 200 ok, got %d", rec.Code)
	}
}

func TestDeleteContact(t *testing.T) {
	t.Cleanup(handlers.storage.Reset)
	handlers.storage.CreateContact(models.Contact{UserId: 1})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", nil)
	ctx := context.WithValue(req.Context(), "user_id", 1)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	handlers.DeleteContact(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("want 200 ok, got %d", rec.Code)
	}
}
