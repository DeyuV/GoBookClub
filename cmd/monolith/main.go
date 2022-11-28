package main

import (
	"GoLang_Backend/app/bookowners"
	"GoLang_Backend/app/books"
	"GoLang_Backend/app/borrowedbooks"
	"GoLang_Backend/app/users"
	"GoLang_Backend/app/waitingbooks"
	"GoLang_Backend/app/wishlist"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	// config
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Database connection
	conn, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// repositories
	usersRepo := users.NewRepository(conn)
	booksRepo := books.NewRepository(conn)
	ownersRepo := bookowners.NewRepository(conn)
	borrowedRepo := borrowedbooks.NewRepository(conn)
	waitingRepo := waitingbooks.NewRepository(conn)
	wishlistRepo := wishlist.NewRepository(conn)

	// services
	booksService := books.NewService(booksRepo)
	usersService := users.NewService(usersRepo)
	bookOwnersService := bookowners.NewService(ownersRepo)
	borrowedService := borrowedbooks.NewService(borrowedRepo)
	waitingService := waitingbooks.NewService(waitingRepo)
	wishlistService := wishlist.NewService(wishlistRepo)

	// transport
	router := mux.NewRouter()

	books.RegisterRoutes(router, booksService)
	bookowners.RegisterRoutes(router, bookOwnersService)
	users.RegisterRoutes(router, usersService)
	borrowedbooks.RegisterRoutes(router, borrowedService)
	waitingbooks.RegisterRoutes(router, waitingService)
	wishlist.RegisterRoutes(router, wishlistService)

	// run server
	err = http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Printf("error listening on port (port already in use?) : %#v", err)
		return
	}
}
