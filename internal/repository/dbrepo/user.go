package dbrepo

import (
	"database/sql"
	"time"

	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/models"
	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) repository.UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (ur *UserRepo) GetUser(id int) (models.User, error) {
	var user models.User

	query := `
	SELECT id, email, password, created_at, updated_at 
	FROM users 
	WHERE id = $1`

	err := ur.db.QueryRow(query, id).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (ur *UserRepo) GetAllUser() ([]models.User, error) {
	users := []models.User{}

	query := `SELECT id, email, created_at, updated_at FROM users`

	rows, err := ur.db.Query(query)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return users, err
	}

	return users, nil
}

func (ur *UserRepo) UpdateUser(u models.User) error {

	var newPassword []byte

	if u.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		newPassword = hashed
	} else {
		newPassword = nil
	}

	query := `
	UPDATE users
	SET 
		email = $1,
		password = COALESCE($2, password), 
		updated_at = $3
	WHERE id = $4
	`

	_, err := ur.db.Exec(query,
		u.Email,
		newPassword,
		time.Now(),
		u.Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepo) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := ur.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
