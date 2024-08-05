package models

type User struct {
	Id int `json:"id"`
	Auth
}

type Contact struct {
	Id           int    `json:"id"`
	UserId       int    `json:"user_id"`
	Nama         string `json:"nama"`
	NomorTelepon string `json:"nomor_telepon"`
	Email        string `json:"email"`
	Alamat       string `json:"alamat"`
}

type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}
