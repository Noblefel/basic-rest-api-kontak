package handlers

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ContactHandlers struct {
}

func NewContactHandlers(db *sql.DB) *ContactHandlers {
	return &ContactHandlers{}
}

func (h *ContactHandlers) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	SendJSON(w, r, http.StatusOK, Response{
		Message: "This is contact no. " + id,
	})
}
