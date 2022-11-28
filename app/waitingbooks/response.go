package waitingbooks

import "time"

type WaitingBookResponse struct {
	UserID     uint      `gorm:"primary_key;auto_increment:false;column:user_id"`
	BookID     uint      `gorm:"primary_key;auto_increment:false;column:book_id"`
	DateOfWait time.Time `gorm:"not_null"`
}
