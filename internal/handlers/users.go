package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/forms"
	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/models"
	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/repository"
	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/repository/dbrepo"
	u "github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/utils"
)

type UserHandlers struct {
	repo repository.UserRepo
}

func NewUserHandlers(db *sql.DB) *UserHandlers {
	return &UserHandlers{
		repo: dbrepo.NewUserRepo(db),
	}
}

func NewTestUserHandlers() *UserHandlers {
	return &UserHandlers{
		repo: dbrepo.NewMockUserRepo(),
	}
}

func (h *UserHandlers) All(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.GetAllUser()
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		u.Message(w, http.StatusInternalServerError, "Error when retrieving all users")
		return
	}

	u.JSON(w, http.StatusOK, u.Response{
		Message: "Users retrieved succesfully",
		Data:    users,
	})
}

func (h *UserHandlers) Get(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	user.Password = ""

	u.JSON(w, http.StatusOK, u.Response{
		Message: "User retrieved succesfully",
		Data:    user,
	})
}

func (h *UserHandlers) Update(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	if err := r.ParseForm(); err != nil {
		u.Message(w, http.StatusBadRequest, "Error parsing form")
		return
	}

	form := forms.New(r.PostForm)
	form.Required("email")
	form.Email("email")
	form.StringMinLength("password", 8)

	if !form.ValidOrErr(w, r) {
		return
	}

	newData := models.User{
		Id:       user.Id,
		Email:    form.Get("email"),
		Password: form.Get("password"),
	}

	if err := h.repo.UpdateUser(newData); err != nil {
		u.Message(w, http.StatusInternalServerError, "Error unable to update user")
		return
	}

	u.Message(w, http.StatusOK, "User updated")
}

func (h *UserHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	if err := h.repo.DeleteUser(user.Id); err != nil {
		u.Message(w, http.StatusInternalServerError, "Error unable to delete user")
		return
	}

	u.Message(w, http.StatusOK, "User deleted")
}
