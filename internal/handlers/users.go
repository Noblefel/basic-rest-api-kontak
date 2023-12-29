package handlers

import (
	"database/sql"
	"errors"
	"net/http"

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

func (h *UserHandlers) All(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.GetAllUser()
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		u.SendJSON(w, r, http.StatusInternalServerError, u.Response{
			Message: "Error when retrieving all users",
		})
		return
	}

	u.SendJSON(w, r, http.StatusOK, u.Response{
		Message: "Users retrieved succesfully",
		Data:    users,
	})
}

func (h *UserHandlers) Get(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	u.SendJSON(w, r, http.StatusOK, u.Response{
		Message: "User retrieved succesfully",
		Data:    user,
	})
}

func (h *UserHandlers) Update(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	err := r.ParseForm()
	if err != nil {
		u.SendJSON(w, r, http.StatusBadRequest, u.Response{
			Message: "Error parsing form",
		})
		return
	}

	newData := models.User{
		Id:       user.Id,
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}

	err = h.repo.UpdateUser(newData)
	if err != nil {
		u.SendJSON(w, r, http.StatusInternalServerError, u.Response{
			Message: "Error unable to update user",
		})
		return
	}

	u.SendJSON(w, r, http.StatusOK, u.Response{
		Message: "User updated",
	})
}

func (h *UserHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	err := h.repo.DeleteUser(user.Id)
	if err != nil {
		u.SendJSON(w, r, http.StatusInternalServerError, u.Response{
			Message: "Error unable to delete user",
		})
		return
	}

	u.SendJSON(w, r, http.StatusOK, u.Response{
		Message: "User deleted",
	})
}
