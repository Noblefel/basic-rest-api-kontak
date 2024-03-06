package dbrepo

import (
	"database/sql"
	"errors"

	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/models"
	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/repository"
)

type mockUserRepo struct{}

func NewMockUserRepo() repository.UserRepo {
	return &mockUserRepo{}
}

func (ur *mockUserRepo) GetUser(id int) (models.User, error) {
	var user models.User

	if id > 9999 {
		return user, sql.ErrNoRows
	}

	if id <= 0 {
		return user, errors.New("unexpected error")
	}

	return user, nil
}

func (ur *mockUserRepo) GetAllUser() ([]models.User, error) {
	var users []models.User

	return users, nil
}

func (ur *mockUserRepo) UpdateUser(u models.User) error {
	if u.Id <= 0 {
		return errors.New("unexpected error")
	}

	return nil
}

func (ur *mockUserRepo) DeleteUser(id int) error {
	if id <= 0 {
		return errors.New("unexpeceted error")
	}

	return nil
}
