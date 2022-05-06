package users

import (
	"context"

	"bitbucket.org/faceit/app"
)

// Errors
var (
	ErrorNotfound         = app.BusinessError("not found")
	ErrEmailAlreadyExists = app.BusinessError("email already exists")
)

type users interface {
	Get(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, userID string) error
	List(ctx context.Context, filter Filter) ([]User, error)
}

type Service struct {
	Users users
}

type UserParams struct {
	FirstName string
	LastName  string
	Nickname  string
	Password  string
	Email     string
	Country   string
}

func (s *Service) Get(ctx context.Context, email string) (*User, error) {
	return s.Users.Get(ctx, email)
}

func (s *Service) Create(ctx context.Context, user *User) error {
	return s.Users.Create(ctx, user)
}

func (s *Service) Update(ctx context.Context, user *User) error {
	return s.Users.Update(ctx, user)
}

func (s *Service) Delete(ctx context.Context, userID string) error {
	return s.Users.Delete(ctx, userID)
}

func (s *Service) List(ctx context.Context, filter Filter) ([]User, error) {
	return s.Users.List(ctx, filter)
}
