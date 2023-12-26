package handlers

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserHandlers struct {
}

func NewUserHandlers(db *sql.DB) *UserHandlers {
	return &UserHandlers{}
}

func (h *UserHandlers) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	SendJSON(w, r, http.StatusOK, Response{
		Message: "This is user no. " + id,
	})
}
