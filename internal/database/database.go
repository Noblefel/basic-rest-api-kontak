package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const maxOpenConns = 10
const maxIdleConns = 5
const maxLifetime = 5 * time.Minute

func Connect() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)

	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(maxOpenConns)
	conn.SetMaxIdleConns(maxIdleConns)
	conn.SetConnMaxLifetime(maxLifetime)

	return conn, nil
}
