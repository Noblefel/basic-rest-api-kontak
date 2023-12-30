package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors,omitempty"`
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
