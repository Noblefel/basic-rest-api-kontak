package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string              `json:"message,omitempty"`
	Data    any                 `json:"data,omitempty"`
	Errors  map[string][]string `json:"errors,omitempty"`
}

func JSON(w http.ResponseWriter, code int, res Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}

// Message is a small wrapper to JSON
func Message(w http.ResponseWriter, code int, msg string) {
	JSON(w, code, Response{Message: msg})
}

// MapString is a tiny helper to append validation error messages
type MapString map[string][]string

func (ms MapString) Add(key, msg string) { ms[key] = append(ms[key], msg) }
