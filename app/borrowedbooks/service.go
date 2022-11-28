package borrowedbooks

import (
	"context"
	"github.com/pkg/errors"
)

type Repository interface {
	Create(ctx context.Context, userId, bookOwnerId, borrowedPeriod uint) (*BorrowResponse, error)
	GetBorrowById(ctx context.Context, id uint) (*BorrowResponse, error)
	ExtendBorrowPeriod(ctx context.Context, id, period uint) (*BorrowResponse, error)
	GetBorrowsByUserId(ctx context.Context, id uint) ([]BorrowResponse, error)
	NotBorrowedBooks(ctx context.Context) ([]NotBorrowedResponse, error)
	BorrowWithReturn(ctx context.Context, userId uint) ([]BorrowWithReturnResponse, error)
}

func NewService(repo Repository) Service {
	return &serviceImplementation{
		repo: repo,
	}
}

type serviceImplementation struct {
	repo Repository
}

func (s *serviceImplementation) GetNotBorrowedBooks(ctx context.Context) ([]NotBorrowedResponse, error) {
	return s.repo.NotBorrowedBooks(ctx)
}

func (s *serviceImplementation) CreateBorrow(ctx context.Context, userId, bookId, period uint) (*BorrowResponse, error) {
	notBorrowed, err := s.repo.NotBorrowedBooks(ctx)
	if err != nil {
		return nil, err
	}

	for _, book := range notBorrowed {
		if book.BookId == bookId {
			return s.repo.Create(ctx, userId, book.BookOwnerId, period)
		}
	}
	return nil, errors.New("the book cant be borrowed for now")
}

func (s *serviceImplementation) ExtendPeriod(ctx context.Context, id, period uint) (*BorrowResponse, error) {
	return s.repo.ExtendBorrowPeriod(ctx, id, period)
}

func (s *serviceImplementation) GetBorrowedBooks(ctx context.Context, userId uint) ([]BorrowWithReturnResponse, error) {
	return
}
