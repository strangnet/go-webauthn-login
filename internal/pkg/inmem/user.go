package inmem

import (
	"github.com/pkg/errors"
	"github.com/strangnet/go-webauthn-login/internal/pkg/domain"
	"sync"
)

type UserRepository struct {
	users map[string]*domain.User
	mu sync.RWMutex
}

func (u *UserRepository) FindByUsername(username string) (*domain.User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	user, ok := u.users[username]
	if !ok {
		return nil, errors.New("User does not exist")
	}

	return user, nil
}

func (u *UserRepository) Create(user *domain.User) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	uu, _ := u.users[user.Username()]
	if uu != nil {
		return errors.New("User already exist")
	}
	u.users[user.Username()] = user
	return nil
}

func (u *UserRepository) ListAllUsers() []*domain.User {
	u.mu.Lock()
	defer u.mu.Unlock()
	var users []*domain.User
	for _, u := range u.users {
		users = append(users, u)
	}

	return users
}

func NewUserRepository() domain.UserRepository {
	return &UserRepository{
		users: make(map[string]*domain.User),
	}
}
