package datastores

import (
	"errors"
	"sync"
	"time"

	"lysk-battle-record/internal/models"
	"lysk-battle-record/internal/sheet_clients"

	"github.com/sirupsen/logrus"
)

type UserStore interface {
	Get(id string) (models.User, bool)
	Insert(user models.User)
	Update(user models.User) error
}

type InMemoryUserStore struct {
	mu          sync.RWMutex
	users       []models.User
	sheetClient sheet_clients.UserSheetClient
}

func NewInMemoryUserStore(sheetClient sheet_clients.UserSheetClient) *InMemoryUserStore {
	store := &InMemoryUserStore{
		sheetClient: sheetClient,
	}
	go store.autoRefresh()
	return store
}

func (s *InMemoryUserStore) autoRefresh() {
	for {
		s.refresh()
		time.Sleep(5 * time.Minute)
	}
}

func (s *InMemoryUserStore) refresh() {
	data, err := s.sheetClient.FetchAllSheetData()
	if err != nil {
		logrus.Errorf("failed to refresh cache for sheet: %s with error %v", s.sheetClient.GetType(), err)
		return
	}

	s.mu.Lock()
	s.users = data
	s.mu.Unlock()

	logrus.Infof("sheet %s refreshed %d users", s.sheetClient.GetType(), len(s.users))
}

func (s *InMemoryUserStore) Get(id string) (models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, u := range s.users {
		if u.ID == id {
			return u, true
		}
	}
	return models.User{}, false
}

func (s *InMemoryUserStore) Insert(user models.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users = append(s.users, user)
}

func (s *InMemoryUserStore) Update(user models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, u := range s.users {
		if u.ID == user.ID {
			s.users[i] = user
			return nil
		}
	}
	return errors.New("user not found")
}
