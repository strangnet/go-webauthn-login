package user

import (
	"github.com/strangnet/go-webauthn-login/internal/pkg/domain"
)

type Service interface {
	FindByUsername(username string) (*domain.User, error)
	Create(user *domain.User) error
	ListAllUsers() []*domain.User
	}

type service struct {
	users domain.UserRepository
}

func NewService(users domain.UserRepository) Service {
	return &service{
		users: users,
	}
}

func (s *service) FindByUsername(username string) (*domain.User, error) {
	return s.users.FindByUsername(username)
}

func (s *service) Create(user *domain.User) error {
	return s.users.Create(user)
}

func (s *service) ListAllUsers() []*domain.User {
	return s.users.ListAllUsers()
}
