package usecase

import (
	"test-assigment/internal/modules/movies/repo"
	"test-assigment/internal/modules/movies/types"

	"github.com/google/uuid"
)

type Service struct {
	repo repo.Movies
}

func New(r repo.Movies) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) TriggerPanic() {
	s.repo.TriggerPanic()
}

func (s *Service) GetMovies() ([]types.Movie, error) {
	return s.repo.GetMovies()
}

func (s *Service) CreateMovie(movie types.Movie) (uuid.UUID, error) {
	return s.repo.CreateMovie(movie)
}

func (s *Service) DeleteMovie(ID uuid.UUID) error {
	return s.repo.DeleteMovie(ID)
}

func (s *Service) GetMovieById(ID uuid.UUID) (types.Movie, error) {
	return s.repo.GetMovieById(ID)
}

func (s *Service) GetMovieByAuthor(ID uuid.UUID) ([]types.Movie, error) {
	return s.repo.GetMovieByAuthor(ID)
}

func (s *Service) GetAuthorById(AuthorID uuid.UUID) (types.Author, error) {
	return s.repo.GetAuthorById(AuthorID)
}

func (s *Service) GetAuthors() ([]types.Author, error) {
	return s.repo.GetAuthors()
}

func (s *Service) CreateAuthor(author types.Author) (uuid.UUID, error) {
	return s.repo.CreateAuthor(author)
}

func (s *Service) DeleteAuthor(ID uuid.UUID) error {
	return s.repo.DeleteAuthor(ID)
}

func (s *Service) GetAuthorWithMovies(AuthorID uuid.UUID) ([]types.Author, error) {
	return s.repo.GetAuthorWithMovies(AuthorID)
}
