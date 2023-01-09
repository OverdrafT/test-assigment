package repo

import (
	"test-assigment/internal/modules/movies/types"

	"gorm.io/gorm"
)

type Movies interface {
	GetMovies() (movies *gorm.DB, err error)
	CreateMovie(movie types.Movie) (DBid string, err error)
	DeleteMovie(movieID string) error
}
