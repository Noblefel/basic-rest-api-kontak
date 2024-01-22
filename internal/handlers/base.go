package handlers

import (
	"net/http"

	u "github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/utils"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	u.SendJSON(w, http.StatusNotFound, u.Response{
		Message: "Not Found",
	})
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	u.SendJSON(w, http.StatusMethodNotAllowed, u.Response{
		Message: "Method Not Allowed",
	})
}
