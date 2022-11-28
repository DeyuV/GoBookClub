package waitingbooks

import "context"

type Repository interface {
	Create(ctx context.Context, userId, bookId uint) (*WaitingBookResponse, error)
}

func NewService(repo Repository) Service {
	return &serviceImplementation{
		repo: repo,
	}
}

type serviceImplementation struct {
	repo Repository
}

func (s *serviceImplementation) CreateWaitingBook(ctx context.Context, userId, bookId uint) (*WaitingBookResponse, error) {
	return s.repo.Create(ctx, userId, bookId)
}
