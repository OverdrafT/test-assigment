package orientdb

import "test-assigment/internal/modules/movies/types"

type orientdbGetAuthorsResponse struct {
	Result []types.Author `json:"result"`
}

type orientdbCreateAuthorResponse struct {
	types.Author
	Class string `json:"class"`
}

type orientdbGetMoviesResponse struct {
	Result []types.Movie `json:"result"`
}

type orientdbGetMovieByID struct {
	Result []types.Movie `json:"result"`
}

type orientdbCreateMovieRequest struct {
	types.Movie
	Class string `json:"@class"`
}
