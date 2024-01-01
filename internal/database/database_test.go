package database

import (
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestConnect(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_PORT", "5432")
		os.Setenv("DB_NAME", "managemen_kontak")
		os.Setenv("DB_USER", "postgres")
		os.Setenv("DB_PASSWORD", "")

		conn, err := Connect()
		if err != nil {
			t.Errorf("Connect() expected no error but got %v", err)
		}
		conn.Close()

		if conn == nil {
			t.Errorf("Connect() wants *sql.DB but got nil")
		}
	})

	t.Run("Fail", func(t *testing.T) {
		os.Setenv("DB_HOST", "x")
		os.Setenv("DB_PORT", "x")
		os.Setenv("DB_NAME", "x")
		os.Setenv("DB_USER", "x")
		os.Setenv("DB_PASSWORD", "x")

		_, err := Connect()
		if err == nil {
			t.Errorf("Connect() expected error but got none")
		}
	})
}
