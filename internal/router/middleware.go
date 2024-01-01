package router

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/models"
	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/repository"
	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/repository/dbrepo"
	u "github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/utils"
	"github.com/go-chi/chi/v5"
)

type Middleware struct {
	user    repository.UserRepo
	contact repository.ContactRepo
}

func NewMiddleware(db *sql.DB) *Middleware {
	return &Middleware{
		user:    dbrepo.NewUserRepo(db),
		contact: dbrepo.NewContactRepo(db),
	}
}

func NewTestMiddleware() *Middleware {
	return &Middleware{
		user:    dbrepo.NewTestUserRepo(),
		contact: dbrepo.NewTestContactRepo(),
	}
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			u.SendJSON(w, r, http.StatusUnauthorized, u.Response{
				Message: "Unauthorized",
			})
			return
		}

		userId, level, err := u.VerifyJWT(tokenString)
		if err != nil {
			u.SendJSON(w, r, http.StatusUnauthorized, u.Response{
				Message: "Unauthorized",
			})
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", int(userId))
		ctx = context.WithValue(ctx, "level", int(level))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) UserGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("user_id").(int)
		userLevel := r.Context().Value("level").(int)
		userIdRoute, err := strconv.Atoi(chi.URLParam(r, "user_id"))
		if err != nil {
			u.SendJSON(w, r, http.StatusBadRequest, u.Response{
				Message: "Invalid id",
			})
			return
		}

		if userId != userIdRoute && userLevel != models.ROLE_ADMIN {
			u.SendJSON(w, r, http.StatusUnauthorized, u.Response{
				Message: "Unauthorized - Sorry you have no permission to do that",
			})
			return
		}

		user, err := m.user.GetUser(userIdRoute)
		if err != nil {
			if errors.Is(sql.ErrNoRows, err) {
				u.SendJSON(w, r, http.StatusNotFound, u.Response{
					Message: "User not found",
				})
				return
			}

			u.SendJSON(w, r, http.StatusInternalServerError, u.Response{
				Message: "Error when retrieving a user",
			})
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) ContactGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("user_id").(int)
		userLevel := r.Context().Value("level").(int)
		contactId, err := strconv.Atoi(chi.URLParam(r, "contact_id"))
		if err != nil {
			u.SendJSON(w, r, http.StatusBadRequest, u.Response{
				Message: "Invalid contact id",
			})
			return
		}

		contact, err := m.contact.GetContact(contactId)
		if err != nil {
			if errors.Is(sql.ErrNoRows, err) {
				u.SendJSON(w, r, http.StatusNotFound, u.Response{
					Message: "Contact not found",
				})
				return
			}

			u.SendJSON(w, r, http.StatusInternalServerError, u.Response{
				Message: "Error when retrieving contact",
			})
			return
		}

		if userId != contact.UserId && userLevel != models.ROLE_ADMIN {
			u.SendJSON(w, r, http.StatusUnauthorized, u.Response{
				Message: "Unauthorized - Sorry, you have no permission to do that",
			})
			return
		}

		ctx := context.WithValue(r.Context(), "contact", contact)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userLevel := r.Context().Value("level").(int)

		if userLevel != models.ROLE_ADMIN {
			u.SendJSON(w, r, http.StatusUnauthorized, u.Response{
				Message: "Unauthorized - Sorry you have no permission to do that",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
