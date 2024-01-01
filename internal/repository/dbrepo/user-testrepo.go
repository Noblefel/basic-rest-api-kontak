package dbrepo

import (
	"database/sql"
	"errors"

	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/models"
)

func (ur *testUserRepo) GetUser(id int) (models.User, error) {
	var user models.User

	if id > 9999 {
		return user, sql.ErrNoRows
	}

	if id <= 0 {
		return user, errors.New("unexpected error")
	}

	return user, nil
}

func (ur *testUserRepo) GetAllUser() ([]models.User, error) {
	var users []models.User

	return users, nil
}

func (ur *testUserRepo) UpdateUser(u models.User) error {
	if u.Id <= 0 {
		return errors.New("unexpected error")
	}

	return nil
}

func (ur *testUserRepo) DeleteUser(id int) error {
	if id <= 0 {
		return errors.New("unexpeceted error")
	}

	return nil
}
