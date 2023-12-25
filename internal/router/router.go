package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type router struct {
}

func NewRouter() *router {
	return &router{}
}

func (r *router) Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	return mux
}
