package books

import (
	"context"
)

type Repository interface {
	GetSuggestedBooks(ctx context.Context, searchTerm string) ([]Book, error)
	CreateBook(ctx context.Context, title, author string, year uint) (*Book, error)
	GetBookByTitleAndAuthor(ctx context.Context, book Book) (*Book, error)
	GetBooksByTitleAndAuthor(ctx context.Context, title, author string) ([]BookData, error)
	FindBooksThatWereBorrowed(ctx context.Context, bookIds []BookData) ([]BookData, error)
}

func NewService(repo Repository) Service {
	return &serviceImplementation{
		repo: repo,
	}
}

type serviceImplementation struct {
	repo Repository
}

func (s *serviceImplementation) PostCreateBook(ctx context.Context, title, author string, year uint) (*Book, error) {
	return s.repo.CreateBook(ctx, title, author, year)
}

func (s *serviceImplementation) GetBooksIDByTitleOrAuthor(ctx context.Context, title, author string) ([]BookData, error) {
	books, err := s.repo.GetBooksByTitleAndAuthor(ctx, title, author)
	if err != nil {
		return nil, err
	}

	return s.repo.FindBooksThatWereBorrowed(ctx, books)
}

func (s *serviceImplementation) GetSuggestedBooks(ctx context.Context, searchTerm string) ([]Book, error) {
	return s.repo.GetSuggestedBooks(ctx, searchTerm)
}
