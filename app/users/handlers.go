package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Service interface {
	CreateUser(ctx context.Context, FirstName, LastName, Email, Password string) (*UserResponse, error)
}

// POSTCreateUser create user account
func POSTCreateUser(svc Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user CreateUserRequest

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			return
		}

		newUser, err := svc.CreateUser(r.Context(), user.FirstName, user.LastName, user.Email, user.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(err.Error())
			if err != nil {
				return
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(newUser)
		if err != nil {
			return
		}
	})
}

func RegisterRoutes(router *mux.Router, svc Service) {
	router.Handle("/user/create", POSTCreateUser(svc)).Methods(http.MethodPost)
}
