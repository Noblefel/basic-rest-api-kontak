package dbrepo

import (
	"database/sql"
	"time"

	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/models"
	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) repository.AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (ar *AuthRepo) Register(u models.User) (int, error) {
	query := `
	INSERT INTO users (email, password, created_at, updated_at) 
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`

	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	var id int

	err = ar.db.QueryRow(query,
		u.Email,
		string(password),
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (ar *AuthRepo) Authenticate(u models.User) (int, error) {
	query := `SELECT id, password FROM users WHERE email = $1`

	var id int
	var password string

	err := ar.db.QueryRow(query, u.Email).Scan(&id, &password)
	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(u.Password))
	if err != nil {
		return 0, err
	}

	return id, nil
}
