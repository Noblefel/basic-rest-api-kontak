package router

import (
	"database/sql"
	"net/http"

	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/handlers"
	"github.com/go-chi/chi/v5"
)

type router struct {
	db      *sql.DB
	auth    *handlers.AuthHandlers
	user    *handlers.UserHandlers
	contact *handlers.ContactHandlers
}

func NewRouter(db *sql.DB) *router {
	return &router{
		db:      db,
		auth:    handlers.NewAuthHandlers(db),
		user:    handlers.NewUserHandlers(db),
		contact: handlers.NewContactHandlers(db),
	}
}

func (r *router) Routes() http.Handler {
	mux := chi.NewRouter()
	middleware := NewMiddleware(r.db)

	mux.Route("/auth", func(mux chi.Router) {
		mux.Post("/register", r.auth.Register)
		mux.Post("/login", r.auth.Login)
	})

	mux.Group(func(mux chi.Router) {
		mux.Use(middleware.Auth)

		mux.Get("/users", r.user.All)
		mux.Route("/users/{user_id}", func(mux chi.Router) {
			mux.Use(middleware.UserGuard)
			mux.Get("/", r.user.Get)
			mux.Post("/update", r.user.Update)
			mux.Post("/delete", r.user.Delete)
			mux.Get("/contacts", r.contact.GetByUser)
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
