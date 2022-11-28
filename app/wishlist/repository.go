package wishlist

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewRepository(db *pgxpool.Pool) Repository {
	return &repositoryImpl{db: db}
}

type repositoryImpl struct {
	db *pgxpool.Pool
}

func (r *repositoryImpl) Create(ctx context.Context, wish Wishlist) (*Wishlist, error) {
	_, err := r.db.Query(ctx, "INSERT INTO wishlists (user_id, book_id) VALUES ($1,$2)", wish.UserID, wish.BookID)
	if err != nil {
		return nil, err
	}

	return &wish, nil
}

func (r *repositoryImpl) GetWishlistByUserId(ctx context.Context, id uint) ([]Wishlist, error) {
	var wishlist []Wishlist

	rows, err := r.db.Query(ctx, "SELECT * FROM wishlists WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		wishlist = append(wishlist, Wishlist{
			UserID: uint(values[0].(int32)),
			BookID: uint(values[1].(int32)),
		})
	}
	return wishlist, nil
}

func (r *repositoryImpl) Delete(ctx context.Context, userId, bookId uint) error {
	_, err := r.db.Query(ctx, "delete from wishlists where user_id = $1 and book_id = $2", userId, bookId)
	if err != nil {
		return err
	}
	return nil
}
