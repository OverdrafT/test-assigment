package transport

import (
	"net/http"
)

type Transport interface {
	GetMovies(w http.ResponseWriter, r *http.Request)
	CreateMovie(w http.ResponseWriter, r *http.Request)
	DeleteMovie(w http.ResponseWriter, r *http.Request)
	GetMovieById(w http.ResponseWriter, r *http.Request)
}
