//go:build integration

package waitingbooks

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

func TestCreateWaiting(t *testing.T) {
	conn := InitConn(t)
	defer func() {
		sqlDB, err := conn.DB()
		if err != nil {
			log.Fatalln(err)
		}
		sqlDB.Close()
	}()

	repo := NewRepository(conn)
	newWait, err := repo.Create(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("error while creating wait : %#v", err)
	}
	t.Logf("new wait : %#v", newWait)
}
