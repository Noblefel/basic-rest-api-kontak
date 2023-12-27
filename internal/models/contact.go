package models

import "time"

type Contact struct {
	Id           int       `json:"id"`
	UserId       int       `json:"user_id"`
	Nama         string    `json:"nama"`
	NomorTelepon string    `json:"nomor_telepon"`
	Email        string    `json:"email"`
	Alamat       string    `json:"alamat"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	User         *User     `json:"user,omitempty"`
}
