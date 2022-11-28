package books

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"strconv"
	"time"
)

//TODO : give up on gorm, use pgx pool, so we can use context cancellation

func NewRepository(db *pgxpool.Pool) Repository {
	return &repositoryImpl{db: db}
}

type repositoryImpl struct {
	db *pgxpool.Pool
}

func (r *repositoryImpl) GetSuggestedBooks(ctx context.Context, searchTerm string) ([]Book, error) {
	var result []Book

	rows, err := r.db.Query(ctx, "SELECT * FROM books WHERE title LIKE  '%' || $1 || '%' OR author LIKE '%' || $1 || '%'", searchTerm)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		values, err2 := rows.Values()
		if err2 != nil {
			return nil, err2
		}
		book := Book{
			ID:     uint(values[0].(int32)),
			Title:  values[1].(string),
			Author: values[2].(string),
			Year:   uint(values[3].(int32)),
		}

		result = append(result, book)
	}
	return result, nil
}

func (r *repositoryImpl) CreateBook(ctx context.Context, title, author string, year uint) (*Book, error) {
	createdBook := Book{
		ID:     0,
		Title:  title,
		Author: author,
		Year:   year,
	}
	err := r.db.QueryRow(ctx, "INSERT INTO books (title, author, published_year) VALUES ($1,$2,$3) RETURNING book_id", title, author, year).Scan(&createdBook.ID)
	if err != nil {
		return nil, err
	}
	return &createdBook, nil
}

// GetBookByTitleAndAuthor return the id of book
func (r *repositoryImpl) GetBookByTitleAndAuthor(ctx context.Context, book Book) (*Book, error) {
	var searchedBook Book

	err := r.db.QueryRow(ctx, "SELECT * FROM books WHERE title = $1 AND author = $2", book.Title, book.Author).Scan(&searchedBook)
	if err != nil {
		return nil, err
	}
	return &searchedBook, nil
}

func (r *repositoryImpl) GetBooksByTitleAndAuthor(ctx context.Context, title, author string) ([]BookData, error) {
	var result []BookData

	err := r.db.QueryRow(ctx, "SELECT * FROM books WHERE title = $1 AND author = $2", title, author).Scan(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repositoryImpl) FindBooksThatWereBorrowed(ctx context.Context, books []BookData) ([]BookData, error) {
	// TODO : refactor without gorm , so we can use on JOIN to rule them all
	bookIds := make([]uint, 0)
	for _, book := range books {
		bookIds = append(bookIds, book.ID)
	}

	type BookList struct {
		BookID uint `gorm:"primary_key;auto_increment:false;column:book_id"`
		RentID uint
	}
	var bookList []BookList
	err := r.db.QueryRow(ctx, "SELECT * FROM bookowners WHERE book_id = $1", bookIds).Scan(&bookList)
	if err != nil {
		return nil, err
	}

	var rentIds []uint
	for _, book := range bookList {
		if book.RentID != 0 {
			rentIds = append(rentIds, book.RentID)
		}
	}

	type RentList struct {
		BookID     uint      `gorm:"primary_key;auto_increment:false;column:book_id"`
		Period     string    `gorm:"not_null"`
		DateOfRent time.Time `gorm:"not_null"`
	}
	var rents []RentList
	r.db.QueryRow(ctx, "SELECT * FROM borrowedbooks where bookowner_id = $1", rentIds).Scan(&rents)

	result := make([]BookData, len(books))
	for i := range books {
		result[i] = books[i]
	}

	for _, rent := range rents {
		num, err2 := strconv.Atoi(string(rent.Period[0]))
		if err2 != nil {
			return nil, err2
		}

		date := rent.DateOfRent.AddDate(0, 0, num*7)
		for i, book := range result {
			if book.ID == rent.BookID {
				result[i].DateOfReturn = &date
				break
			}
		}
	}

	return result, nil
}
