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
		repo: dbrepo.NewTestUserRepo(),
	}
}

func (h *UserHandlers) All(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.GetAllUser()
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		u.SendJSON(w, http.StatusInternalServerError, u.Response{
			Message: "Error when retrieving all users",
		})
		return
	}

	u.SendJSON(w, http.StatusOK, u.Response{
		Message: "Users retrieved succesfully",
		Data:    users,
	})
}

func (h *UserHandlers) Get(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	u.SendJSON(w, http.StatusOK, u.Response{
		Message: "User retrieved succesfully",
		Data:    user,
	})
}

func (h *UserHandlers) Update(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	err := r.ParseForm()
	if err != nil {
		u.SendJSON(w, http.StatusBadRequest, u.Response{
			Message: "Error parsing form",
		})
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

	err = h.repo.UpdateUser(newData)
	if err != nil {
		u.SendJSON(w, http.StatusInternalServerError, u.Response{
			Message: "Error unable to update user",
		})
		return
	}

	u.SendJSON(w, http.StatusOK, u.Response{
		Message: "User updated",
	})
}

func (h *UserHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	err := h.repo.DeleteUser(user.Id)
	if err != nil {
		u.SendJSON(w, http.StatusInternalServerError, u.Response{
			Message: "Error unable to delete user",
		})
		return
	}

	u.SendJSON(w, http.StatusOK, u.Response{
		Message: "User deleted",
	})
}
