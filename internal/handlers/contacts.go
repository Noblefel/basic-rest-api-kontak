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

type ContactHandlers struct {
	repo repository.ContactRepo
}

func NewContactHandlers(db *sql.DB) *ContactHandlers {
	return &ContactHandlers{
		repo: dbrepo.NewContactRepo(db),
	}
}

func NewTestContactHandlers() *ContactHandlers {
	return &ContactHandlers{
		repo: dbrepo.NewTestContactRepo(),
	}
}

func (h *ContactHandlers) All(w http.ResponseWriter, r *http.Request) {
	contacts, err := h.repo.GetAllContact()
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		u.SendJSON(w, http.StatusInternalServerError, u.Response{
			Message: "Error retrieving all contact",
		})
		return
	}

	u.SendJSON(w, http.StatusOK, u.Response{
		Message: "Contacts retrieved",
		Data:    contacts,
	})
}

func (h *ContactHandlers) Get(w http.ResponseWriter, r *http.Request) {
	contact := r.Context().Value("contact").(models.Contact)

	u.SendJSON(w, http.StatusOK, u.Response{
		Message: "Contact retrieved",
		Data:    contact,
	})
}

func (h *ContactHandlers) GetByUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	contacts, err := h.repo.GetUserContact(user.Id)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		u.SendJSON(w, http.StatusInternalServerError, u.Response{
			Message: "Error retrieving user's contacts",
		})
		return
	}

	u.SendJSON(w, http.StatusOK, u.Response{
		Message: "User's contacts retrieved",
		Data:    contacts,
	})
}

func (h *ContactHandlers) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		u.SendJSON(w, http.StatusBadRequest, u.Response{
			Message: "Error parsing form",
		})
		return
	}

	form := forms.New(r.PostForm)
	form.Required("nama")
	form.Email("email")

	if !form.ValidOrErr(w, r) {
		return
	}

	userId := r.Context().Value("user_id").(int)

	contact := models.Contact{
		UserId:       userId,
		Nama:         form.Get("nama"),
		NomorTelepon: form.Get("nomor_telepon"),
		Email:        form.Get("email"),
		Alamat:       form.Get("alamat"),
	}

	id, err := h.repo.CreateContact(contact)
	if err != nil {
		u.SendJSON(w, http.StatusInternalServerError, u.Response{
			Message: "Error unable to create contact",
		})
		return
	}

	u.SendJSON(w, http.StatusCreated, u.Response{
		Message: "Contact created",
		Data:    id,
	})
}

func (h *ContactHandlers) Update(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		u.SendJSON(w, http.StatusBadRequest, u.Response{
			Message: "Error parsing form",
		})
		return
	}

	form := forms.New(r.PostForm)
	form.Required("nama")
	form.Email("email")

	if !form.ValidOrErr(w, r) {
		return
	}

	contact := r.Context().Value("contact").(models.Contact)
	contact = models.Contact{
		Id:           contact.Id,
		Nama:         form.Get("nama"),
		NomorTelepon: form.Get("nomor_telepon"),
		Email:        form.Get("email"),
		Alamat:       form.Get("alamat"),
	}

	err = h.repo.UpdateContact(contact)
	if err != nil {
		u.SendJSON(w, http.StatusInternalServerError, u.Response{
			Message: "Error unable to update contact",
		})
		return
	}

	u.SendJSON(w, http.StatusOK, u.Response{
		Message: "Contact updated",
	})
}

func (h *ContactHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	contact := r.Context().Value("contact").(models.Contact)

	err := h.repo.DeleteContact(contact.Id)
	if err != nil {
		u.SendJSON(w, http.StatusInternalServerError, u.Response{
			Message: "Error unable to delete contact",
		})
		return
	}

	u.SendJSON(w, http.StatusOK, u.Response{
		Message: "Contact deleted",
	})
}
