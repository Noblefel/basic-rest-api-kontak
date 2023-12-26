package router

import (
	"net/http"

	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/handlers"
	"github.com/go-chi/chi/v5"
)

type router struct {
	user    *handlers.UserHandlers
	contact *handlers.ContactHandlers
}

func NewRouter(
	user *handlers.UserHandlers,
	contact *handlers.ContactHandlers,
) *router {
	return &router{
		user:    user,
		contact: contact,
	}
}

func (r *router) Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Route("/users", func(mux chi.Router) {
		mux.Get("/{id}", r.user.Get)
	})

	mux.Route("/contacts", func(mux chi.Router) {
		mux.Get("/{id}", r.contact.Get)
	})

	mux.NotFound(handlers.NotFound)
	mux.MethodNotAllowed(handlers.MethodNotAllowed)

	return mux
}
