package forms

import (
	"fmt"
	"net/http"
	"net/mail"
	"net/url"
	"strings"

	u "github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/utils"
)

type Form struct {
	url.Values
	Errors errors
}

type errors map[string][]string

func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

func New(data url.Values) *Form {
	return &Form{
		data,
		make(errors),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)

		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "Field cannot be empty")
		}
	}
}

func (f *Form) StringMinLength(field string, n int) {
	value := strings.TrimSpace(f.Get(field))

	// This is for optional form fields
	if value == "" {
		return
	}

	if len(value) < n {
		f.Errors.Add(field, fmt.Sprintf("Field must be atleast %d characters", n))
	}
}

func (f *Form) Email(field string) {
	// This is for optional form fields
	if strings.TrimSpace(f.Get(field)) == "" {
		return
	}

	_, err := mail.ParseAddress(f.Get(field))
	if err != nil {
		f.Errors.Add(field, "Email is not valid")
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) ValidOrErr(w http.ResponseWriter, r *http.Request) bool {
	if !f.Valid() {
		u.SendJSON(w, http.StatusBadRequest, u.Response{
			Message: "Some fields are invalid",
			Errors:  f.Errors,
		})
		return false
	}

	return true
}
