package handlers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendJSON(w http.ResponseWriter, r *http.Request, code int, res Response) {
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Error encoding JSON Response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonBytes)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	SendJSON(w, r, http.StatusNotFound, Response{
		Message: "Not Found",
	})
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	SendJSON(w, r, http.StatusMethodNotAllowed, Response{
		Message: "Method Not Allowed",
	})
}
