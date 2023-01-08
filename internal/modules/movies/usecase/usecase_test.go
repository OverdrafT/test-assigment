package usecase

import (
	"errors"
	"reflect"
	"test-assigment/internal/modules/movies/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMovies(t *testing.T) {
	s := New(repoMock{})

	got, err := s.GetMovies()

	want := []types.Movie{
		{
			MovieYear: 2000,
			MovieName: "help",
		},
	}

	// if got != want {
	// 	t.Errorf("got %q, wanted %q", got, want)
	// }
	assert.Equal(t, true, reflect.DeepEqual(got, want))

	assert.Nil(t, err)
}

func TestGetMoviesError(t *testing.T) {
	s := New(repoMock{
		err: errors.New("Failed!"),
	})

	_, err := s.GetMovies()

	assert.Error(t, err)
}

type repoMock struct {
	err error
}

func (r repoMock) GetMovies() (movies []types.Movie, err error) {
	return []types.Movie{
		{
			MovieYear: 2000,
			MovieName: "help",
		},
	}, r.err
}

func (r repoMock) CreateMovie(movie types.Movie) (DBid string, err error) {
	panic("Implement me!")
}

func (r repoMock) DeleteMovie(movieID string) error {
	panic("Implement me!")
}
