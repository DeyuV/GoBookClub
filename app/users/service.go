package users

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	VerifyUser(ctx context.Context, userId uint) error
	GetUserById(ctx context.Context, userid uint) (*UserResponse, error)
	Create(ctx context.Context, FirstName, LastName, Email, PasswordHash string, Enabled bool, Role string) (*UserResponse, error)
}

func NewService(repo Repository) Service {
	return &serviceImplementation{
		repo: repo,
	}
}

type serviceImplementation struct {
	repo Repository
}

func (s *serviceImplementation) CreateUser(ctx context.Context, FirstName, LastName, Email, Password string) (*UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), 14)
	if err != nil {
		return nil, err
	}

	createdUser, err := s.repo.Create(ctx, FirstName, LastName, Email, string(hashedPassword), true, "USER")
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
