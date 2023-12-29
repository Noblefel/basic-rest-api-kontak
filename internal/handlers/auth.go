package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/models"
	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/repository"
	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/repository/dbrepo"
	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/utils"
	u "github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandlers struct {
	repo repository.AuthRepo
}

func NewAuthHandlers(db *sql.DB) *AuthHandlers {
	return &AuthHandlers{
		repo: dbrepo.NewAuthRepo(db),
	}
}

func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		u.SendJSON(w, r, http.StatusBadRequest, u.Response{
			Message: "Error parsing form",
		})
		return
	}

	user := models.User{
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}

	_, err = h.repo.Register(user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			u.SendJSON(w, r, http.StatusConflict, u.Response{
				Message: "Email already in use",
			})
			return
		}

		u.SendJSON(w, r, http.StatusInternalServerError, u.Response{
			Message: "Error unable to register user",
		})
		return
	}

	u.SendJSON(w, r, http.StatusCreated, u.Response{
		Message: "User has been created",
	})
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		u.SendJSON(w, r, http.StatusBadRequest, u.Response{
			Message: "Error parsing form",
		})
		return
	}

	user := models.User{
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}

	id, err := h.repo.Authenticate(user)
	if err != nil {
		if errors.Is(bcrypt.ErrMismatchedHashAndPassword, err) || errors.Is(sql.ErrNoRows, err) {
			u.SendJSON(w, r, http.StatusUnauthorized, u.Response{
				Message: "Invalid credentials",
			})
			return
		}

		u.SendJSON(w, r, http.StatusInternalServerError, u.Response{
			Message: "Error when authenticating",
		})
		return
	}

	token, err := utils.GenerateJWT(id)
	if err != nil {
		u.SendJSON(w, r, http.StatusInternalServerError, u.Response{
			Message: "Error when authenticating",
		})
		return
	}

	u.SendJSON(w, r, http.StatusOK, u.Response{
		Message: "Authenticated",
		Data:    token,
	})
}
