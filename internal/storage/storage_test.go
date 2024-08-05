package storage

import (
	"fmt"
	"sync"
	"testing"

	"github.com/Noblefel/baic-rest-api-kontak/internal/models"
)

func TestStorage(t *testing.T) {
	store := New()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			store.Register(fmt.Sprint(i), "")
			wg.Done()
		}()

		go func() {
			store.CreateContact(models.Contact{})
			wg.Done()
		}()
	}

	wg.Wait()

	if len(store.users) != 100 || len(store.contacts) != 100 {
		t.Errorf(
			"expecting 100 inserted data for each, got %d users, %d contacts",
			len(store.users),
			len(store.contacts),
		)
	}

	store.Reset()

	if len(store.users) != 0 || len(store.contacts) != 0 {
		t.Errorf(
			"storage should be empty, counted %d users, %d contacts",
			len(store.users),
			len(store.contacts),
		)
	}
}
