package borrowedbooks

import "time"

type BorrowResponse struct {
	ID             uint
	UserID         uint
	BookOwnerID    uint
	BorrowedPeriod uint
	BorrowedDate   time.Time
}

type BorrowWithReturnResponse struct {
	UserID       uint
	FirstName    string
	LastName     string
	Title        string
	Author       string
	Year         string
	DateOfReturn time.Time
}

type NotBorrowedResponse struct {
	BookId      uint
	BookOwnerId uint
	Title       string
	Author      string
	Year        uint
}
