package main

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Noblefel/baic-rest-api-kontak/internal/utils"
)

func TestServer(t *testing.T) {
	go func() {
		main()
	}()

	time.Sleep(time.Second)
	_, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Error(err)
	}
}

func TestRoutes(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("", "/api/ping", nil)

	routes := routes()
	routes.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("want 200 ok, got %d", rec.Code)
	}

	rec = httptest.NewRecorder()
	req = httptest.NewRequest("", "/404", nil)
	routes.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("want 404 not found, got %d", rec.Code)
	}
}

func TestAuth(t *testing.T) {
	appStorage.Register("", "")

	token, err := utils.GenerateJWT(1)
	if err != nil {
		t.Fatal(err)
	}

	handler := auth(func(w http.ResponseWriter, r *http.Request) {
		id, ok := r.Context().Value("user_id").(int)
		if !ok {
			t.Fatal("type assertion error")
		}

		if id != 1 {
			t.Errorf("want id of 1, got %d", id)
		}
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", token)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("want 200 ok, got %d", rec.Code)
	}
}
