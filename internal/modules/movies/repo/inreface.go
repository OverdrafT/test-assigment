package repo

import "test-assigment/internal/modules/movies/types"

type Movies interface {
	GetMovies() (movies []types.Movie, err error)
	CreateMovie(movie types.Movie) (DBid string, err error)
	DeleteMovie(movieID string) error
}
