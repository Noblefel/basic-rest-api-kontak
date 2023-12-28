package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/models"
	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/repository"
	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
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
		SendJSON(w, r, http.StatusInternalServerError, Response{
			Message: "Error when retrieving all users",
		})
		return
	}

	SendJSON(w, r, http.StatusOK, Response{
		Message: "Users retrieved succesfully",
		Data:    users,
	})
}

func (h *UserHandlers) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		SendJSON(w, r, http.StatusBadRequest, Response{
			Message: "Invalid id",
		})
		return
	}

	user, err := h.repo.GetUser(id)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			SendJSON(w, r, http.StatusNotFound, Response{
				Message: "User not found",
			})
			return
		}

		SendJSON(w, r, http.StatusInternalServerError, Response{
			Message: "Error when retrieving a user",
		})
		return
	}

	SendJSON(w, r, http.StatusOK, Response{
		Message: "User retrieved succesfully",
		Data:    user,
	})
}

func (h *UserHandlers) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		SendJSON(w, r, http.StatusBadRequest, Response{
			Message: "Invalid id",
		})
		return
	}

	_, err = h.repo.GetUser(id)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			SendJSON(w, r, http.StatusNotFound, Response{
				Message: "User not found",
			})
			return
		}

		SendJSON(w, r, http.StatusInternalServerError, Response{
			Message: "Error when retrieving a user",
		})
		return
	}

	err = r.ParseForm()
	if err != nil {
		SendJSON(w, r, http.StatusBadRequest, Response{
			Message: "Error parsing form",
		})
		return
	}

	newData := models.User{
		Id:       id,
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}

	err = h.repo.UpdateUser(newData)
	if err != nil {
		SendJSON(w, r, http.StatusInternalServerError, Response{
			Message: "Error unable to update user",
		})
		return
	}

	SendJSON(w, r, http.StatusOK, Response{
		Message: "User updated",
	})
}

func (h *UserHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		SendJSON(w, r, http.StatusBadRequest, Response{
			Message: "Invalid id",
		})
		return
	}

	_, err = h.repo.GetUser(id)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			SendJSON(w, r, http.StatusNotFound, Response{
				Message: "User not found",
			})
			return
		}

		SendJSON(w, r, http.StatusInternalServerError, Response{
			Message: "Error when retrieving a user",
		})
		return
	}

	err = h.repo.DeleteUser(id)
	if err != nil {
		SendJSON(w, r, http.StatusInternalServerError, Response{
			Message: "Error unable to delete user",
		})
		return
	}

	SendJSON(w, r, http.StatusOK, Response{
		Message: "User deleted",
	})
}
