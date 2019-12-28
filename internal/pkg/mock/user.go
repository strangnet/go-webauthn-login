package mock

import "github.com/strangnet/go-webauthn-login/internal/pkg/domain"

type UserRepository struct {
	CreateFn func(*domain.User) error
	CreateInvoked bool

	FindByUsernameFn      func(string) (*domain.User, error)
	FindByUsernameInvoked bool

	ListAllUsersFn      func() []*domain.User
	ListAllUsersInvoked bool
}

func (u *UserRepository) FindByUsername(username string) (*domain.User, error) {
	u.FindByUsernameInvoked = true
	return u.FindByUsernameFn(username)
}

func (u *UserRepository) Create(user *domain.User) error {
	u.CreateInvoked = true
	return u.CreateFn(user)
}

func (u *UserRepository) ListAllUsers() []*domain.User {
	u.ListAllUsersInvoked = true
	return u.ListAllUsersFn()
}

func NewUserRepository() domain.UserRepository {
	return &UserRepository{}
}
