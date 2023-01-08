package memory

import (
	"test-assigment/internal/modules/movies/types"
)

type Service struct {
	movies []types.Movie
}

func New() *Service {
	return &Service{
		movies: []types.Movie{
			{
				MovieYear: 2000,
				MovieName: "help",
			},
			{
				MovieYear: 2003,
				MovieName: "Hello World",
			},
		},
	}
}

func (s *Service) GetMovies() []types.Movie {
	return s.movies
}
