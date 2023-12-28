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
		SendJSON(w, r, http.StatusBadRequest, Response{
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
			SendJSON(w, r, http.StatusConflict, Response{
				Message: "Email already in use",
			})
			return
		}

		SendJSON(w, r, http.StatusInternalServerError, Response{
			Message: "Error unable to register user",
		})
		return
	}

	SendJSON(w, r, http.StatusCreated, Response{
		Message: "User has been created",
	})
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		SendJSON(w, r, http.StatusBadRequest, Response{
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
			SendJSON(w, r, http.StatusUnauthorized, Response{
				Message: "Invalid credentials",
			})
			return
		}

		SendJSON(w, r, http.StatusInternalServerError, Response{
			Message: "Error when authenticating",
		})
		return
	}

	token, err := utils.GenerateJWT(id)
	if err != nil {
		SendJSON(w, r, http.StatusInternalServerError, Response{
			Message: "Error when authenticating",
		})
		return
	}

	SendJSON(w, r, http.StatusOK, Response{
		Message: "Authenticated",
		Data:    token,
	})
}
