package transport

import (
	"net/http"
)

type Transport interface {
	GetMovies(w http.ResponseWriter, r *http.Request)
	CreateMovie(w http.ResponseWriter, r *http.Request)
	DeleteMovie(w http.ResponseWriter, r *http.Request)
	GetMovieById(w http.ResponseWriter, r *http.Request)

	GetMovieByAuthor(w http.ResponseWriter, r *http.Request)

	GetAuthors(w http.ResponseWriter, r *http.Request)
	CreateAuthor(w http.ResponseWriter, r *http.Request)
	GetAuthorById(w http.ResponseWriter, r *http.Request)
	DeleteAuthor(w http.ResponseWriter, r *http.Request)
	GetAuthorWithMovies(w http.ResponseWriter, r *http.Request)

	TriggerPanic(http.ResponseWriter, *http.Request)
}
