package bookowners

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, userId, bookId uint) (*BookOwner, error)
}

func NewService(repo Repository) Service {
	return &serviceImplementation{
		repo: repo,
	}
}

type serviceImplementation struct {
	repo Repository
}

func (s *serviceImplementation) PostCreateBookOwners(ctx context.Context, UserID, BookID uint) (*BookOwner, error) {
	return s.repo.Create(ctx, UserID, BookID)
}
