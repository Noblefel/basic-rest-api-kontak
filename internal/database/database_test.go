package database

import (
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestConnect(t *testing.T) {
	t.Run("succcess", func(t *testing.T) {
		_, err := Connect("pgx", "host=localhost port=5432 dbname=managemen_kontak user=postgres password=")
		if err != nil {
			t.Errorf("Connect() expected no error but got %v", err)
		}
	})

	t.Run("fail opening connection", func(t *testing.T) {
		_, err := Connect("", "")
		if err == nil {
			t.Errorf("Connect() expected error")
		}
	})

	t.Run("fail pinging", func(t *testing.T) {
		_, err := Connect("pgx", "")
		if err == nil {
			t.Errorf("Connect() expected error")
		}
	})
}
