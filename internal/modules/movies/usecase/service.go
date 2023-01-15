package usecase

import (
	"github.com/google/uuid"
	"test-assigment/internal/modules/movies/repo"
	"test-assigment/internal/modules/movies/types"
)

type Service struct {
	repo repo.Movies
}

func New(r repo.Movies) *Service {
	return &Service{
		repo: r,
	}
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
