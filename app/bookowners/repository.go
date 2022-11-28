package bookowners

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

func (r *repositoryImpl) Create(ctx context.Context, userId, bookId uint) (*BookOwner, error) {
	bookOwner := BookOwner{
		Id:     0,
		UserID: userId,
		BookID: bookId,
	}

	err := r.db.QueryRow(ctx, "INSERT INTO bookowners (user_id, book_id) VALUES ($1,$2) RETURNING bookowner_id", userId, bookId).Scan(&bookOwner.Id)
	if err != nil {
		return nil, err
	}
	return &bookOwner, nil
}
