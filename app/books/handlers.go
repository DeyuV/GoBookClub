package books

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Data Transfer Object = DTO - makes the translation from transport into an object / struct that can be used by the Service

type Service interface {
	PostCreateBook(ctx context.Context, title, author string, year uint) (*Book, error)
	GetSuggestedBooks(ctx context.Context, searchTerm string) ([]Book, error)
	GetBooksIDByTitleOrAuthor(ctx context.Context, title, author string) ([]BookData, error)
}

func POSTCreateBook(svc Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bookRequest BookRequest
		err := json.NewDecoder(r.Body).Decode(&bookRequest)
		if err != nil {
			return
		}

		result, err := svc.PostCreateBook(r.Context(), bookRequest.Title, bookRequest.Author, bookRequest.Year)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(err)
			if err != nil {
				return
			}
		}

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(&result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

// GETSuggestBook suggest a book by title or author containing the string
func GETSuggestBook(svc Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request SuggestRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			// maybe we should say something to the client, like 403 bad request or 500 server error
			return
		}

		result, err := svc.GetSuggestedBooks(r.Context(), request.Term)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(&result)
		if err != nil {
			// maybe we should say something to the client, like 500 server error
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

// GETSearchBookByTitleOrAuthor get all books by matching title or author and see if it is available or when it will be available
func GETSearchBookByTitleOrAuthor(svc Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request BookByAuthorOrTitleRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			// answer with error to client
			return
		}

		result, err := svc.GetBooksIDByTitleOrAuthor(r.Context(), request.Title, request.Author)
		if err != nil {
			// TODO : answer with error to client
			return
		}
		err = json.NewEncoder(w).Encode(&result)
		if err != nil {
			return
		}
	})
}

func RegisterRoutes(router *mux.Router, svc Service) {
	router.Handle("/book/create", POSTCreateBook(svc)).Methods(http.MethodPost)
	router.Handle("/book/suggest", GETSuggestBook(svc)).Methods(http.MethodGet)
	router.Handle("/booklist/search", GETSearchBookByTitleOrAuthor(svc)).Methods(http.MethodGet)
}
