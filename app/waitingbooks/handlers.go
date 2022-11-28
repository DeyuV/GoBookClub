package waitingbooks

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Service interface {
	CreateWaitingBook(ctx context.Context, userId, bookId uint) (*WaitingBookResponse, error)
}

func POSTCreateWaitingBook(svc Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var waitBook WaitBookRequest

		err := json.NewDecoder(r.Body).Decode(&waitBook)
		if err != nil {
			return
		}

		newWaitBook, err := svc.CreateWaitingBook(r.Context(), waitBook.UserID, waitBook.BookID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(err.Error())
			if err != nil {
				return
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(&newWaitBook)
		if err != nil {
			return
		}
	})
}
func RegisterRoutes(router *mux.Router, svc Service) {
	router.Handle("/wait/create", POSTCreateWaitingBook(svc)).Methods(http.MethodPost)
}
