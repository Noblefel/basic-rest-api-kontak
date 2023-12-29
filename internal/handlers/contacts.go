package handlers

import (
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

type ContactHandlers struct {
	repo repository.ContactRepo
}

func NewContactHandlers(db *sql.DB) *ContactHandlers {
	return &ContactHandlers{
		repo: dbrepo.NewContactRepo(db),
	}
}

func (h *ContactHandlers) All(w http.ResponseWriter, r *http.Request) {
	contacts, err := h.repo.GetAllContact()
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		u.SendJSON(w, r, http.StatusInternalServerError, u.Response{
			Message: "Error retrieving all contact",
		})
		return
	}

	u.SendJSON(w, r, http.StatusOK, u.Response{
		Message: "Contacts retrieved",
		Data:    contacts,
	})
}

func (h *ContactHandlers) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "contact_id"))
	if err != nil {
		u.SendJSON(w, r, http.StatusBadRequest, u.Response{
			Message: "Invalid id",
		})
		return
	}

	contact, err := h.repo.GetContactWithUser(id)
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

	u.SendJSON(w, r, http.StatusOK, u.Response{
		Message: "Contact retrieved",
		Data:    contact,
	})
}

func (h *ContactHandlers) GetByUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		u.SendJSON(w, r, http.StatusBadRequest, u.Response{
			Message: "Invalid id",
		})
		return
	}

	contacts, err := h.repo.GetUserContact(userId)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		u.SendJSON(w, r, http.StatusInternalServerError, u.Response{
			Message: "Error retrieving user's contacts",
		})
		return
	}

	u.SendJSON(w, r, http.StatusOK, u.Response{
		Message: "User's contacts retrieved",
		Data:    contacts,
	})
}

func (h *ContactHandlers) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		u.SendJSON(w, r, http.StatusBadRequest, u.Response{
			Message: "Error parsing form",
		})
		return
	}

	contact := models.Contact{
		UserId:       1,
		Nama:         r.Form.Get("nama"),
		NomorTelepon: r.Form.Get("nomor_telepon"),
		Email:        r.Form.Get("email"),
		Alamat:       r.Form.Get("alamat"),
	}

	id, err := h.repo.CreateContact(contact)
	if err != nil {
		u.SendJSON(w, r, http.StatusInternalServerError, u.Response{
			Message: "Error unable to create contact",
		})
		return
	}

	u.SendJSON(w, r, http.StatusCreated, u.Response{
		Message: "Contact created",
		Data:    id,
	})
}

func (h *ContactHandlers) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "contact_id"))
	if err != nil {
		u.SendJSON(w, r, http.StatusBadRequest, u.Response{
			Message: "Invalid id",
		})
		return
	}

	_, err = h.repo.GetContact(id)
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

	err = r.ParseForm()
	if err != nil {
		u.SendJSON(w, r, http.StatusBadRequest, u.Response{
			Message: "Error parsing form",
		})
		return
	}

	contact := models.Contact{
		Id:           id,
		Nama:         r.Form.Get("nama"),
		NomorTelepon: r.Form.Get("nomor_telepon"),
		Email:        r.Form.Get("email"),
		Alamat:       r.Form.Get("alamat"),
	}

	err = h.repo.UpdateContact(contact)
	if err != nil {
		u.SendJSON(w, r, http.StatusInternalServerError, u.Response{
			Message: "Error unable to update contact",
		})
		return
	}

	u.SendJSON(w, r, http.StatusOK, u.Response{
		Message: "Contact updated",
	})
}

func (h *ContactHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "contact_id"))
	if err != nil {
		u.SendJSON(w, r, http.StatusBadRequest, u.Response{
			Message: "Invalid id",
		})
		return
	}

	_, err = h.repo.GetContact(id)
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

	err = h.repo.DeleteContact(id)
	if err != nil {
		u.SendJSON(w, r, http.StatusInternalServerError, u.Response{
			Message: "Error unable to delete contact",
		})
		return
	}

	u.SendJSON(w, r, http.StatusOK, u.Response{
		Message: "Contact deleted",
	})
}
