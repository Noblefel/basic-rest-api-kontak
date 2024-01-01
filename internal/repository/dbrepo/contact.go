package dbrepo

import (
	"database/sql"
	"time"

	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/models"
	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/repository"
)

type ContactRepo struct {
	db *sql.DB
}

type testContactRepo struct{}

func NewContactRepo(db *sql.DB) repository.ContactRepo {
	return &ContactRepo{
		db: db,
	}
}

func NewTestContactRepo() repository.ContactRepo {
	return &testContactRepo{}
}

func (cr *ContactRepo) GetContact(id int) (models.Contact, error) {
	var contact models.Contact

	query := `
	SELECT id, user_id, nama, nomor_telepon, email, alamat, created_at, updated_at 
	FROM contacts
	WHERE id = $1
	`

	err := cr.db.QueryRow(query, id).Scan(
		&contact.Id,
		&contact.UserId,
		&contact.Nama,
		&contact.NomorTelepon,
		&contact.Email,
		&contact.Alamat,
		&contact.CreatedAt,
		&contact.UpdatedAt,
	)

	if err != nil {
		return contact, err
	}

	return contact, nil
}

func (cr *ContactRepo) GetContactWithUser(id int) (models.Contact, error) {
	var contact models.Contact
	contact.User = &models.User{}

	query := `
	SELECT c.id, c.user_id, c.nama, c.nomor_telepon, c.email, c.alamat, c.created_at, c.updated_at, u.id, u.email
	FROM contacts c
	LEFT JOIN users u ON (c.user_id = u.id)
	WHERE c.id = $1
	`

	err := cr.db.QueryRow(query, id).Scan(
		&contact.Id,
		&contact.UserId,
		&contact.Nama,
		&contact.NomorTelepon,
		&contact.Email,
		&contact.Alamat,
		&contact.CreatedAt,
		&contact.UpdatedAt,
		&contact.User.Id,
		&contact.User.Email,
	)

	if err != nil {
		return contact, err
	}

	return contact, nil
}

func (cr *ContactRepo) GetAllContact() ([]models.Contact, error) {
	contacts := []models.Contact{}

	query := `
	SELECT id, user_id, nama, nomor_telepon, email, alamat, created_at, updated_at 
	FROM contacts
	`

	rows, err := cr.db.Query(query)
	if err != nil {
		return contacts, err
	}

	for rows.Next() {
		var contact models.Contact

		err := rows.Scan(
			&contact.Id,
			&contact.UserId,
			&contact.Nama,
			&contact.NomorTelepon,
			&contact.Email,
			&contact.Alamat,
			&contact.CreatedAt,
			&contact.UpdatedAt,
		)

		if err != nil {
			return contacts, nil
		}

		contacts = append(contacts, contact)
	}

	if err = rows.Err(); err != nil {
		return contacts, err
	}

	return contacts, nil
}

func (cr *ContactRepo) GetUserContact(user_id int) ([]models.Contact, error) {
	contacts := []models.Contact{}

	query := `
	SELECT id, user_id, nama, nomor_telepon, email, alamat, created_at, updated_at 
	FROM contacts
	WHERE user_id = $1
	`
	rows, err := cr.db.Query(query, user_id)
	if err != nil {
		return contacts, err
	}

	for rows.Next() {
		var contact models.Contact

		err := rows.Scan(
			&contact.Id,
			&contact.UserId,
			&contact.Nama,
			&contact.NomorTelepon,
			&contact.Email,
			&contact.Alamat,
			&contact.CreatedAt,
			&contact.UpdatedAt,
		)

		if err != nil {
			return contacts, nil
		}

		contacts = append(contacts, contact)
	}

	if err = rows.Err(); err != nil {
		return contacts, err
	}

	return contacts, nil
}

func (cr *ContactRepo) CreateContact(c models.Contact) (int, error) {
	query := `
	INSERT INTO contacts (user_id, nama, nomor_telepon, email, alamat, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`

	var id int

	err := cr.db.QueryRow(query,
		c.UserId,
		c.Nama,
		c.NomorTelepon,
		c.Email,
		c.Alamat,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (cr *ContactRepo) UpdateContact(c models.Contact) error {
	query := `
	UPDATE contacts
	SET
		nama = $1,
		nomor_telepon = $2,
		email = $3,
		alamat = $4,
		updated_at = $5
	WHERE id = $6
	`

	_, err := cr.db.Exec(query,
		c.Nama,
		c.NomorTelepon,
		c.Email,
		c.Alamat,
		time.Now(),
		c.Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (cr *ContactRepo) DeleteContact(id int) error {
	query := `DELETE FROM contacts WHERE id = $1`

	_, err := cr.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
