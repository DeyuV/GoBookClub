//go:build integration

package borrowedbooks

/*
import (
	"context"
	"github.com/pkg/errors"
	"testing"
	"time"
)

type mockedRepository struct {
	borrows map[uint]*BorrowResponse
}

func (m *mockedRepository) Create(ctx context.Context, userId, bookId, period uint) (*BorrowResponse, error) {

	for _, borrow := range m.borrows {
		if borrow.UserID == userId && borrow.BookID == bookId {
			return nil, errors.New("the borrow already exists")
		}
	}

	borrowCount := uint(len(m.borrows))

	result := BorrowResponse{
		ID:         borrowCount,
		UserID:     userId,
		BookID:     bookId,
		Period:     period,
		DateOfRent: time.Now(),
	}
	m.borrows[borrowCount] = &result
	return &result, nil
}

func newMockedRepository() Repository {
	result := mockedRepository{
		borrows: make(map[uint]*BorrowResponse),
	}

	result.Create(context.Background(), 1, 1, 2)
	result.Create(context.Background(), 1, 2, 1)
	return &result
}

func TestCreateServiceBorrow(t *testing.T) {
	repo := newMockedRepository()
	svc := NewService(repo)
	newBorrow, err := svc.CreateBorrow(context.Background(), 1, 4, 4)
	if err != nil {
		t.Fatalf("error creating borrow : %#v", err)
	}
	t.Logf("borrow was created : %#v", newBorrow)
}
*/
