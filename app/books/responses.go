package books

import (
	"time"
)

type Book struct {
	ID     uint   `gorm:"primary_key;column:book_id"`
	Title  string `gorm:"not_null"`
	Author string `gorm:"not_null"`
	Year   uint   `gorm:"not_null"`
}

type BookData struct {
	ID           uint       `json:"id"`                     // id of the book
	Title        string     `json:"title" gorm:"not_null"`  //
	Author       string     `json:"author" gorm:"not_null"` //
	Year         uint       `json:"year" gorm:"not_null"`   //
	DateOfReturn *time.Time `json:"dateOfReturn,omitempty"` // if set, book was borrowed
}
