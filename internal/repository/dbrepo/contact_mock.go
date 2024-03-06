package dbrepo

import (
	"database/sql"
	"errors"

	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/models"
	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/repository"
)

type mockContactRepo struct{}

func NewMockContactRepo() repository.ContactRepo {
	return &mockContactRepo{}
}

func (cr *mockContactRepo) GetContact(id int) (models.Contact, error) {
	contact := models.Contact{
		Id:     1,
		UserId: 1,
	}

	if id > 9999 {
		return contact, sql.ErrNoRows
	}

	if id <= 0 {
		return contact, errors.New("unexpected error")
	}

	return contact, nil
}

func (cr *mockContactRepo) GetContactWithUser(id int) (models.Contact, error) {
	var contact models.Contact

	return contact, nil
}

func (cr *mockContactRepo) GetAllContact() ([]models.Contact, error) {
	var contacts []models.Contact

	return contacts, nil
}

func (cr *mockContactRepo) GetUserContact(user_id int) ([]models.Contact, error) {
	var contacts []models.Contact

	if user_id <= 0 {
		return contacts, errors.New("unexpected error")
	}

	return contacts, nil
}

func (cr *mockContactRepo) CreateContact(c models.Contact) (int, error) {
	if c.UserId <= 0 {
		return 0, errors.New("unexpected error")
	}

	return 1, nil
}

func (cr *mockContactRepo) UpdateContact(c models.Contact) error {
	if c.Id <= 0 {
		return errors.New("unexpected error")
	}

	return nil
}

func (cr *mockContactRepo) DeleteContact(id int) error {
	if id <= 0 {
		return errors.New("unexpected error")
	}

	return nil
}
