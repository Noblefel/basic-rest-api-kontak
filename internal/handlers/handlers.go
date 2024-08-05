package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/Noblefel/baic-rest-api-kontak/internal/models"
	"github.com/Noblefel/baic-rest-api-kontak/internal/storage"
	u "github.com/Noblefel/baic-rest-api-kontak/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type Handlers struct{ storage *storage.Storage }

func New(s *storage.Storage) *Handlers {
	return &Handlers{s}
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ping"))
}

func (app *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	var body models.Auth

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err)
		u.Message(w, http.StatusBadRequest, "tidak bisa proses json")
		return
	}

	valError := make(u.MapString)

	if _, err := mail.ParseAddress(body.Email); err != nil {
		valError.Add("email", "email tidak benar")
	}

	if len(body.Password) < 8 {
		valError.Add("password", "password minimal 8 karakter")
	}

	if len(body.Password) > 50 {
		valError.Add("password", "password kepanjangan")
	}

	if len(valError) != 0 {
		u.JSON(w, http.StatusBadRequest, u.Response{
			Message: "kesalahan input",
			Errors:  valError,
		})
		return
	}

	_, err = app.storage.GetUserByEmail(body.Email)
	if err == nil {
		u.Message(w, http.StatusConflict, storage.ErrUsedEmail.Error())
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		log.Println(err)
		u.Message(w, http.StatusUnauthorized, "ada masalah di server kami")
		return
	}

	app.storage.Register(body.Email, string(hash))
	u.Message(w, http.StatusCreated, "sukses register")
}

func (app *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var body models.Auth

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err)
		u.Message(w, http.StatusBadRequest, "tidak bisa proses json")
		return
	}
	valError := make(u.MapString)

	if _, err := mail.ParseAddress(body.Email); err != nil {
		valError.Add("email", "email tidak benar")
	}

	if len(body.Password) < 8 {
		valError.Add("password", "password minmal 8 karakter")
	}

	if len(body.Password) > 50 {
		valError.Add("password", "password kepanjangan")
	}

	if len(valError) != 0 {
		u.JSON(w, http.StatusBadRequest, u.Response{
			Message: "kesalahan input",
			Errors:  valError,
		})
		return
	}

	user, err := app.storage.GetUserByEmail(body.Email)
	if err != nil {
		u.Message(w, http.StatusUnauthorized, "email atau password salah")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		log.Println(err)
		u.Message(w, http.StatusUnauthorized, "email atau password salah")
		return
	}

	token, err := u.GenerateJWT(user.Id)
	if err != nil {
		log.Println(err)
		u.Message(w, http.StatusInternalServerError, "ada masalah di server kami")
		return
	}

	u.JSON(w, http.StatusOK, u.Response{Data: token})
}

func (app *Handlers) GetAccount(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id").(int)

	user, err := app.storage.GetUser(id)
	if err != nil {
		u.Message(w, http.StatusNotFound, err.Error())
		return
	}

	user.Password = ""
	u.JSON(w, http.StatusOK, u.Response{Data: user})
}

func (app *Handlers) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id").(int)

	user, err := app.storage.GetUser(id)
	if err != nil {
		u.Message(w, http.StatusNotFound, err.Error())
		return
	}

	var body models.Auth
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err)
		u.Message(w, http.StatusBadRequest, "tidak bisa proses json")
		return
	}

	valError := make(u.MapString)

	if _, err := mail.ParseAddress(body.Email); err != nil {
		valError.Add("email", "email tidak benar")
	}

	if len(valError) != 0 {
		u.JSON(w, http.StatusBadRequest, u.Response{
			Message: "kesalahan input",
			Errors:  valError,
		})
		return
	}

	user.Email = body.Email
	app.storage.UpdateUser(user)
	u.Message(w, http.StatusOK, "user telah di update")
}

func (app *Handlers) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id").(int)
	app.storage.DeleteUser(id)
	u.Message(w, http.StatusOK, "akun ini berhasil dihapus")
}

func (app *Handlers) GetAccountContacts(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id").(int)
	contacts := app.storage.GetUserContacts(id)
	u.JSON(w, http.StatusOK, u.Response{Data: contacts})
}

func (app *Handlers) GetContact(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println(err)
		u.Message(w, http.StatusBadRequest, "parameter id tidak benar")
		return
	}

	contact, err := app.storage.GetContact(id)
	if err != nil {
		u.Message(w, http.StatusNotFound, err.Error())
		return
	}

	authId := r.Context().Value("user_id").(int)

	if contact.UserId != authId {
		u.Message(w, http.StatusUnauthorized, "kamu tidak punya hak akses untuk melihat ini")
		return
	}

	u.JSON(w, http.StatusOK, u.Response{Data: contact})
}

func (app *Handlers) CreateContact(w http.ResponseWriter, r *http.Request) {
	var body models.Contact

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err)
		u.Message(w, http.StatusBadRequest, "tidak bisa proses json")
		return
	}

	body.Nama = strings.TrimSpace(body.Nama)
	valError := make(u.MapString)

	if body.Nama == "" {
		valError.Add("nama", "nama masih kosong")
	}

	if _, err := mail.ParseAddress(body.Email); err != nil {
		valError.Add("email", "email tidak benar")
	}

	if len(valError) != 0 {
		u.JSON(w, http.StatusBadRequest, u.Response{
			Message: "kesalahan input",
			Errors:  valError,
		})
		return
	}

	body.UserId = r.Context().Value("user_id").(int)
	id := app.storage.CreateContact(body)

	u.JSON(w, http.StatusCreated, u.Response{
		Message: "kontak telah dibuat",
		Data:    id,
	})
}

func (app *Handlers) UpdateContact(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println(err)
		u.Message(w, http.StatusBadRequest, "parameter id tidak benar")
		return
	}

	contact, err := app.storage.GetContact(id)
	if err != nil {
		u.Message(w, http.StatusNotFound, err.Error())
		return
	}

	if r.Context().Value("user_id").(int) != contact.UserId {
		u.Message(w, http.StatusUnauthorized, "kamu tidak punya hak akses")
		return
	}

	var body models.Contact
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err)
		u.Message(w, http.StatusBadRequest, "tidak bisa proses json")
		return
	}

	body.Nama = strings.TrimSpace(body.Nama)
	valError := make(u.MapString)

	if body.Nama == "" {
		valError.Add("nama", "nama masih kosong")
	}

	if _, err := mail.ParseAddress(body.Email); err != nil {
		valError.Add("email", "email tidak benar")
	}

	if len(valError) != 0 {
		u.JSON(w, http.StatusBadRequest, u.Response{
			Message: "kesalahan input",
			Errors:  valError,
		})
		return
	}

	body.UserId = contact.UserId
	body.Id = contact.Id

	app.storage.UpdateContact(body)
	u.Message(w, http.StatusOK, "kontak telah di update")
}

func (app *Handlers) DeleteContact(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println(err)
		u.Message(w, http.StatusBadRequest, "parameter id tidak benar")
		return
	}

	contact, err := app.storage.GetContact(id)
	if err != nil {
		u.Message(w, http.StatusNotFound, err.Error())
		return
	}

	authId := r.Context().Value("user_id").(int)

	if contact.UserId != authId {
		u.Message(w, http.StatusUnauthorized, "kamu tidak punya hak akses untuk menghapus ini")
		return
	}

	app.storage.DeleteContact(id)
	u.Message(w, http.StatusOK, "kontak telah dihapus")
}
