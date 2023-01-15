package repo

import (
	"github.com/google/uuid"
	"test-assigment/internal/modules/movies/types"
)

type Movies interface {
	GetMovies() (movies []types.Movie, err error)
	CreateMovie(movie types.Movie) (uuid.UUID, error)
	DeleteMovie(ID uuid.UUID) error
	GetMovieById(ID uuid.UUID) (types.Movie, error)
}
