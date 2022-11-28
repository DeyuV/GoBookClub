package waitingbooks

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

func NewRepository(db *pgxpool.Pool) Repository {
	return &repositoryImpl{db: db}
}

type repositoryImpl struct {
	db *pgxpool.Pool
}

func (r *repositoryImpl) Create(ctx context.Context, userId, bookId uint) (*WaitingBookResponse, error) {
	waitingBook := WaitingBookResponse{
		UserID:     userId,
		BookID:     bookId,
		DateOfWait: time.Now(),
	}

	_, err := r.db.Query(ctx, "INSERT INTO waitingbooks (user_id, book_id, waiting_date) VALUES ($1,$2,$3)", userId, bookId, time.Now())
	if err != nil {
		return nil, err
	}

	return &waitingBook, nil
}
