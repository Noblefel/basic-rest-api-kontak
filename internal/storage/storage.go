package storage

import (
	"errors"
	"sync"

	"github.com/Noblefel/baic-rest-api-kontak/internal/models"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrUsedEmail = errors.New("email already used")
)

type Storage struct {
	users map[int]models.User
	rwU   sync.RWMutex
	seqU  int

	contacts map[int]models.Contact
	rwC      sync.RWMutex
	seqC     int
}

func New() *Storage {
	return &Storage{
		users:    make(map[int]models.User),
		contacts: make(map[int]models.Contact),
		rwU:      sync.RWMutex{},
		rwC:      sync.RWMutex{},
	}
}

func (s *Storage) Reset() {
	s.rwU.Lock()
	s.rwC.Lock()
	defer s.rwU.Unlock()
	defer s.rwC.Unlock()

	s.users = make(map[int]models.User)
	s.contacts = make(map[int]models.Contact)

	s.seqU = 0
	s.seqC = 0
}

func (s *Storage) GetUser(id int) (models.User, error) {
	s.rwU.Lock()
	defer s.rwU.Unlock()

	user, ok := s.users[id]
	if !ok {
		return models.User{}, ErrNotFound
	}

	return user, nil
}

func (s *Storage) GetUserByEmail(email string) (*models.User, error) {
	s.rwU.Lock()
	defer s.rwU.Unlock()

	for _, user := range s.users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, ErrNotFound
}

func (s *Storage) Register(email, password string) int {
	var user models.User
	user.Email = email
	user.Password = password

	s.rwU.Lock()
	defer s.rwU.Unlock()

	s.seqU++
	user.Id = s.seqU
	s.users[user.Id] = user

	return user.Id
}

func (s *Storage) UpdateUser(u models.User) {
	s.rwU.Lock()
	defer s.rwU.Unlock()
	s.users[u.Id] = u
}

func (s *Storage) DeleteUser(id int) {
	s.rwU.Lock()
	defer s.rwU.Unlock()
	delete(s.users, id)
}

func (s *Storage) GetContact(id int) (*models.Contact, error) {
	s.rwC.Lock()
	defer s.rwC.Unlock()

	contact, ok := s.contacts[id]
	if !ok {
		return nil, ErrNotFound
	}

	return &contact, nil
}

func (s *Storage) GetUserContacts(userId int) []models.Contact {
	var contacts []models.Contact

	for _, c := range s.contacts {
		if c.UserId == userId {
			contacts = append(contacts, c)
		}
	}

	return contacts
}

func (s *Storage) CreateContact(c models.Contact) int {
	s.rwC.Lock()
	defer s.rwC.Unlock()

	s.seqC++
	c.Id = s.seqC
	s.contacts[c.Id] = c

	return c.Id
}

func (s *Storage) UpdateContact(c models.Contact) {
	s.rwC.Lock()
	defer s.rwC.Unlock()
	s.contacts[c.Id] = c
}

func (s *Storage) DeleteContact(id int) {
	s.rwC.Lock()
	defer s.rwC.Unlock()
	delete(s.contacts, id)
}
