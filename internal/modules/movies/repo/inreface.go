package repo

import (
	"test-assigment/internal/modules/movies/types"

	"github.com/google/uuid"
)

type Movies interface {
	GetMovies() (movies []types.Movie, err error)
	CreateMovie(movie types.Movie) (uuid.UUID, error)
	DeleteMovie(ID uuid.UUID) error
	GetMovieById(ID uuid.UUID) (types.Movie, error)
	GetMovieByAuthor(ID uuid.UUID) ([]types.Movie, error)

	GetAuthors() ([]types.Author, error)
	CreateAuthor(author types.Author) (uuid.UUID, error)
	GetAuthorById(AuthorID uuid.UUID) (types.Author, error)
	DeleteAuthor(ID uuid.UUID) error
	GetAuthorWithMovies(AuthorID uuid.UUID) ([]types.Author, error)

	TriggerPanic()
}
