package wishlist

import (
	"context"
	"github.com/pkg/errors"
)

type Repository interface {
	Create(ctx context.Context, wish Wishlist) (*Wishlist, error)
	GetWishlistByUserId(ctx context.Context, id uint) ([]Wishlist, error)
	Delete(ctx context.Context, userId, bookId uint) error
}

func NewService(repo Repository) Service {
	return &serviceImplementation{
		repo: repo,
	}
}

type serviceImplementation struct {
	repo Repository
}

func (s *serviceImplementation) CreateWish(ctx context.Context, wish Wishlist) (*Wishlist, error) {
	return s.repo.Create(ctx, wish)
}

func (s *serviceImplementation) GetWishByUserId(ctx context.Context, id uint) ([]Wishlist, error) {
	wishList, err := s.repo.GetWishlistByUserId(ctx, id)
	if wishList == nil {
		return nil, errors.New("user doesn't have any books on wishlist")
	}
	return wishList, err
}

func (s *serviceImplementation) DeleteWish(ctx context.Context, userId, bookId uint) error {
	return s.repo.Delete(ctx, userId, bookId)
}
