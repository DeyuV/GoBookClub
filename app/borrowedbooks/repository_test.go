//go:build integration

package borrowedbooks

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

func InitConn(t *testing.T) *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	host := os.Getenv("HOST")
	dbPortStr := os.Getenv("DBPORT")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		t.Fatalf("error reading port : %#v", err.Error())
	}
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")
	// Database connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%d sslmode=disable", host, user, dbName, password, dbPort)
	t.Logf("DB URI : %s", dbURI)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	// Opening connection to database
	conn, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN:                  dbURI,
				PreferSimpleProtocol: true,
			},
		), &gorm.Config{
			Logger: newLogger,
		},
	)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to database!")
	}
	return conn
}

func TestCreateBorrow(t *testing.T) {
	conn := InitConn(t)
	defer func() {
		sqlDB, err := conn.DB()
		if err != nil {
			log.Fatalln(err)
		}
		sqlDB.Close()
	}()

	repo := NewRepository(conn)
	newBorrow, err := repo.Create(context.Background(), 77, 10, 3)
	if err != nil {
		t.Fatalf("error while creating borrow : %#v", err)
	}
	t.Logf("new borrow : %#v", newBorrow)
}

func TestGetBorrowById(t *testing.T) {
	conn := InitConn(t)
	defer func() {
		sqlDB, err := conn.DB()
		if err != nil {
			log.Fatalln(err)
		}
		sqlDB.Close()
	}()

	repo := NewRepository(conn)
	borrow, err := repo.GetBorrowById(context.Background(), 31)
	if err != nil {
		t.Fatalf("error while getting borrow by id: %#v", err)
	}
	t.Logf("got borrow : %#v", borrow)
}

func TestExtendBorrowPeriod(t *testing.T) {
	conn := InitConn(t)
	defer func() {
		sqlDB, err := conn.DB()
		if err != nil {
			log.Fatalln(err)
		}
		sqlDB.Close()
	}()

	repo := NewRepository(conn)
	err := repo.ExtendBorrowPeriod(context.Background(), 31, 2)
	if err != nil {
		t.Fatalf("error while extending period: %#v", err)
	}
	t.Logf("period extended")
}

func TestGetBorrowsByUserId(t *testing.T) {
	conn := InitConn(t)
	defer func() {
		sqlDB, err := conn.DB()
		if err != nil {
			log.Fatalln(err)
		}
		sqlDB.Close()
	}()

	repo := NewRepository(conn)
	borrowedBooks, err := repo.GetBorrowsByUserId(context.Background(), 52)
	if err != nil {
		t.Fatalf("error while getting borrowed books by a user id: %#v", err)
	}
	t.Logf("got borrowed books: %#v", borrowedBooks)
}
