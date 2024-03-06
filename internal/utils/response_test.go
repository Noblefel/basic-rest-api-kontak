package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponse(t *testing.T) {
	var tests = []struct {
		name               string
		data               interface{}
		inputStatusCode    int
		expectedStatusCode int
	}{
		{"success", nil, http.StatusOK, http.StatusOK},
		{"success 2", nil, http.StatusBadRequest, http.StatusBadRequest},
		{"error marshalling json", make(chan int), http.StatusOK, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			JSON(w, tt.inputStatusCode, Response{Data: tt.data})

			if w.Code != tt.expectedStatusCode {
				t.Errorf("want %d, got %d", tt.expectedStatusCode, w.Code)
			}
		})
	}
}
