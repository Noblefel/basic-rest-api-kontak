package handlers

import (
	"net/http"

	u "github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/utils"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	u.Message(w, http.StatusNotFound, "Not Found")
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	u.Message(w, http.StatusMethodNotAllowed, "Method Not Allowed")
}
