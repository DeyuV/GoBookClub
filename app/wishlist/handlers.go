package wishlist

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Service interface {
	CreateWish(ctx context.Context, wish Wishlist) (*Wishlist, error)
	GetWishByUserId(ctx context.Context, id uint) ([]Wishlist, error)
	DeleteWish(ctx context.Context, userId, bookId uint) error
}

func POSTCreateWish(svc Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var wish Wishlist

		err := json.NewDecoder(r.Body).Decode(&wish)
		if err != nil {
			return
		}

		newWish, err := svc.CreateWish(r.Context(), wish)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(err.Error())
			if err != nil {
				return
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(&newWish)
		if err != nil {
			return
		}

	})
}

func GETWishlist(svc Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user UserIdRequest

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			return
		}

		wishList, err := svc.GetWishByUserId(r.Context(), user.UserId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(err.Error())
			if err != nil {
				return
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(wishList)
		if err != nil {
			return
		}
	})
}

func DELETEWish(svc Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var wish Wishlist

		err := json.NewDecoder(r.Body).Decode(&wish)
		if err != nil {
			return
		}

		err = svc.DeleteWish(r.Context(), wish.UserID, wish.BookID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(err.Error())
			if err != nil {
				return
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func RegisterRoutes(router *mux.Router, svc Service) {
	router.Handle("/wishlist/create", POSTCreateWish(svc)).Methods(http.MethodPost)
	router.Handle("/wishlist/get", GETWishlist(svc)).Methods(http.MethodGet)
	router.Handle("/wishlist/delete", DELETEWish(svc)).Methods(http.MethodDelete)
}
