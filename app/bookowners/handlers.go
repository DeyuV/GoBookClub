package bookowners

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Service interface {
	PostCreateBookOwners(ctx context.Context, UserID, BookID uint) (*BookOwner, error)
}

// POSTCreateBookOwner add book for a user in book owners
func POSTCreateBookOwner(svc Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request BookOwnerRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			return
		}

		result, err := svc.PostCreateBookOwners(r.Context(), request.UserID, request.BookID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			if err != nil {
				return
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(&result)
		if err != nil {
			return
		}

	})
}

func RegisterRoutes(router *mux.Router, svc Service) {
	router.Handle("/bookowner/create", POSTCreateBookOwner(svc)).Methods(http.MethodPost)
}
