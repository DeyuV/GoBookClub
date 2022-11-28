package borrowedbooks

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Service interface {
	CreateBorrow(ctx context.Context, userId, bookId, period uint) (*BorrowResponse, error)
	ExtendPeriod(ctx context.Context, id, period uint) (*BorrowResponse, error)
	GetBorrowedBooks(ctx context.Context, userId uint) ([]BorrowedWithReturnResponse, error)
	GetNotBorrowedBooks(ctx context.Context) ([]NotBorrowedResponse, error)
}

func POSTCreateBorrow(svc Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var borrow BorrowRequest

		err := json.NewDecoder(r.Body).Decode(&borrow)
		if err != nil {
			return
		}

		newBorrow, err := svc.CreateBorrow(r.Context(), borrow.UserID, borrow.BookID, borrow.Period)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(err.Error())
			if err != nil {
				return
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(newBorrow)
		if err != nil {
			return
		}
	})
}

func GETNotBorrowedBooks(svc Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result, err := svc.GetNotBorrowedBooks(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(err.Error())
			if err != nil {
				return
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(&result)
		if err != nil {
			return
		}
	})
}

func PUTExtendPeriod(svc Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var extendBorrowRequest ExtendBorrowRequest

		err := json.NewDecoder(r.Body).Decode(&extendBorrowRequest)
		if err != nil {
			return
		}

		extendedBorrow, err := svc.ExtendPeriod(r.Context(), extendBorrowRequest.Id, extendBorrowRequest.Period)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(err.Error())
			if err != nil {
				return
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(extendedBorrow)
		if err != nil {
			return
		}
	})
}

func GETBorrowedBooksWithReturnDate(svc Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user UserIdRequest

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			return
		}

		result, err := svc.GetBorrowedBooks(r.Context(), user.Id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(err.Error())
			if err != nil {
				return
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(&result)
		if err != nil {
			return
		}
	})
}

func RegisterRoutes(router *mux.Router, svc Service) {
	router.Handle("/borrow/create", POSTCreateBorrow(svc)).Methods(http.MethodPost)
	router.Handle("/borrow/not", GETNotBorrowedBooks(svc)).Methods(http.MethodGet)
	router.Handle("/borrow/extend", PUTExtendPeriod(svc)).Methods(http.MethodPut)
	router.Handle("/borrow/return", GETBorrowedBooksWithReturnDate(svc)).Methods(http.MethodGet)
}
