package usecase

import (
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

func (s *Service) CreateMovie(movie types.Movie) (DBid string, err error) {
	return s.repo.CreateMovie(movie)
}

func (s *Service) DeleteMovie(movieID string) error {
	return s.repo.DeleteMovie(movieID)
}
