//go:build integration
// +build integration

package users

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestCreateUser(t *testing.T) {
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

	defer func() {
		sqlDB, err := conn.DB()
		if err != nil {
			log.Fatalln(err)
		}
		sqlDB.Close()
	}()

	repo := NewRepository(conn)
	newUser, err := repo.Create(context.Background(), "Bogdan", "Dinu", "badu@badu.ro", "123456", true, "ADMIN")
	if err != nil {
		t.Fatalf("error while creating user : %#v", err)
	}
	t.Logf("new user : %#v", newUser)
}
