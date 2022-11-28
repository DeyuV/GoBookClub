package borrowedbooks

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

func (r *repositoryImpl) Create(ctx context.Context, userId, bookOwnerId, borrowedPeriod uint) (*BorrowResponse, error) {
	newBorrow := BorrowResponse{
		ID:             0,
		UserID:         userId,
		BookOwnerID:    bookOwnerId,
		BorrowedPeriod: borrowedPeriod,
		BorrowedDate:   time.Now(),
	}

	err := r.db.QueryRow(ctx, `INSERT INTO borrowedbooks (user_id, bookowner_id, borrowed_period, borrowed_date)
								   VALUES ($1,$2,$3,$4) RETURNING borrowed_id`, userId, bookOwnerId, borrowedPeriod, time.Now()).Scan(&newBorrow.ID)
	if err != nil {
		return nil, err
	}

	return &newBorrow, nil
}

func (r *repositoryImpl) NotBorrowedBooks(ctx context.Context) ([]NotBorrowedResponse, error) {
	var notBorrowedBooks []NotBorrowedResponse

	rows, err := r.db.Query(ctx, `select bb.borrowed_id, bo.bookowner_id, bo.book_id,b.title,b.author,b.published_year from borrowedbooks as bb 
									  full outer join bookowners as bo on bb.bookowner_id = bo.bookowner_id 
    								  join books as b on bo.book_id = b.book_id`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		values, err2 := rows.Values()
		if err2 != nil {
			return nil, err2
		}

		if values[0] == nil {
			notBorrowedBooks = append(notBorrowedBooks, NotBorrowedResponse{
				BookOwnerId: uint(values[1].(int32)),
				BookId:      uint(values[2].(int32)),
				Title:       values[3].(string),
				Author:      values[4].(string),
				Year:        uint(values[5].(int32)),
			})
		}
	}
	return notBorrowedBooks, nil
}

func (r *repositoryImpl) GetBorrowById(ctx context.Context, id uint) (*BorrowResponse, error) {
	var borrow BorrowResponse

	err := r.db.QueryRow(ctx, "SELECT * FROM borrowedbooks WHERE bookowner_id = $1", id).Scan(&borrow)
	if err != nil {
		return nil, err
	}
	return &borrow, nil
}

func (r *repositoryImpl) ExtendBorrowPeriod(ctx context.Context, id, period uint) (*BorrowResponse, error) {
	var borrow BorrowResponse

	err := r.db.QueryRow(ctx, "UPDATE borrowedbooks SET borrowed_period = borrowed_period + $1 WHERE borrowed_id = $2 RETURNING *",
		period, id).Scan(&borrow.ID, &borrow.UserID, &borrow.BookOwnerID, &borrow.BorrowedPeriod, &borrow.BorrowedDate)

	if err != nil {
		return nil, err
	}
	return &borrow, nil
}

func (r *repositoryImpl) BorrowWithReturn(ctx context.Context, userId uint) ([]BorrowWithReturnResponse, error) {
	rows, err := r.db.Query(ctx, `select bo.user_id as owner_user_id, bb.user_id as borrow_user_id,
									  u.first_name, u.last_name, bb.borrowed_date, bb.borrowed_period,
									  b.title, b.author, b.published_year from borrowedbooks as bb
									  join users as u on u.user_id = bb.user_id
									  join bookowners as bo on bb.bookowner_id = bo.bookowner_id
 									  join books as b on bo.book_id = b.book_id`)
	if err != nil {
		return nil, err
	}

	var result []BorrowWithReturnResponse
	for rows.Next() {
		values, err2 := rows.Values()
		if err2 != nil {
			return nil, err2
		}

		if values[0] == userId {
			result = append(result, BorrowWithReturnResponse{
				UserID:       uint(values[1].(int32)),
				FirstName:    values[2].(string),
				LastName:     values[3].(string),
				Title:        values[6].(string),
				Author:       values[7].(string),
				Year:         values[8].(string),
				DateOfReturn: values[4].(time.Time).AddDate(0, 0, int(values[5].(int32)*7)),
			})
		}
	}
	return result, nil
}

func (r *repositoryImpl) GetBorrowsByUserId(ctx context.Context, id uint) ([]BorrowResponse, error) {
	var borrowed []BorrowResponse
	err := r.db.QueryRow(ctx, "SELECT * FROM borrowedbooks WHERE user_id = $1", id).Scan(&borrowed)
	if err != nil {
		return nil, err
	}

	return borrowed, nil
}
