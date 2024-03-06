package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/forms"
	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/models"
	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/repository"
	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/repository/dbrepo"
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

func NewTestAuthHandlers() *AuthHandlers {
	return &AuthHandlers{
		repo: dbrepo.NewMockAuthRepo(),
	}
}

func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		u.Message(w, http.StatusBadRequest, "Error parsing form")
		return
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.Email("email")
	form.StringMinLength("password", 8)

	if !form.ValidOrErr(w, r) {
		return
	}

	user := models.User{
		Email:    form.Get("email"),
		Password: form.Get("password"),
	}

	_, err := h.repo.Register(user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			u.Message(w, http.StatusConflict, "Email already in use")
			return
		}

		u.Message(w, http.StatusInternalServerError, "Error unable to register user")
		return
	}

	u.Message(w, http.StatusCreated, "User has been created")
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		u.Message(w, http.StatusBadRequest, "Error parsing form")
		return
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.Email("email")
	form.StringMinLength("password", 8)

	if !form.ValidOrErr(w, r) {
		return
	}

	user := models.User{
		Email:    form.Get("email"),
		Password: form.Get("password"),
	}

	id, level, err := h.repo.Authenticate(user)
	if err != nil {
		if errors.Is(bcrypt.ErrMismatchedHashAndPassword, err) || errors.Is(sql.ErrNoRows, err) {
			u.Message(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		u.Message(w, http.StatusInternalServerError, "Error when authenticating")
		return
	}

	token, err := u.GenerateJWT(id, level)
	if err != nil {
		u.Message(w, http.StatusInternalServerError, "Error when authenticating")
		return
	}

	u.JSON(w, http.StatusOK, u.Response{
		Message: "Authenticated",
		Data:    token,
	})
}
