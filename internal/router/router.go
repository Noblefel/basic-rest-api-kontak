package router

import (
	"net/http"

	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/handlers"
	"github.com/go-chi/chi/v5"
)

type router struct {
	auth    *handlers.AuthHandlers
	user    *handlers.UserHandlers
	contact *handlers.ContactHandlers
}

func NewRouter(
	auth *handlers.AuthHandlers,
	user *handlers.UserHandlers,
	contact *handlers.ContactHandlers,
) *router {
	return &router{
		auth:    auth,
		user:    user,
		contact: contact,
	}
}

func (r *router) Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Route("/auth", func(mux chi.Router) {
		mux.Post("/register", r.auth.Register)
		mux.Post("/login", r.auth.Login)
	})

	mux.Group(func(mux chi.Router) {
		mux.Use(Auth)

		mux.Route("/users", func(mux chi.Router) {
			mux.Get("/", r.user.All)
			mux.Get("/{user_id}", r.user.Get)
			mux.Post("/{user_id}/update", r.user.Update)
			mux.Post("/{user_id}/delete", r.user.Delete)
			mux.Get("/{user_id}/contacts", r.contact.GetByUser)
		})

		mux.Route("/contacts", func(mux chi.Router) {
			mux.Get("/", r.contact.All)
			mux.Post("/create", r.contact.Create)
			mux.Get("/{contact_id}", r.contact.Get)
			mux.Post("/{contact_id}/update", r.contact.Update)
			mux.Post("/{contact_id}/delete", r.contact.Delete)
		})
	})

	mux.NotFound(handlers.NotFound)
	mux.MethodNotAllowed(handlers.MethodNotAllowed)

	return mux
}
