package router

import (
	"database/sql"
	"reflect"
	"testing"
)

func TestNewRouter(t *testing.T) {
	var db *sql.DB
	router := NewRouter(db)

	typeString := reflect.TypeOf(router).String()
	if typeString != "*router.router" {
		t.Error("NewRouter() did not get the correct type, wanted *router.router")
	}
}

func TestRouter_Routes(t *testing.T) {
	var db *sql.DB
	router := NewRouter(db)

	mux := router.Routes()

	typeString := reflect.TypeOf(mux).String()
	if typeString != "*chi.Mux" {
		t.Error("router Routes() did not get the correct type, wanted *chi.Mux")
	}
}
