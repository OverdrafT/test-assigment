package transport

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"test-assigment/internal/modules/movies/types"
	"test-assigment/internal/modules/movies/usecase"

	"github.com/gorilla/mux"
)

const YearOfTheFirstEverMovie = 1895

type service struct {
	uc *usecase.Service
}

func New(us *usecase.Service) *service {
	return &service{
		uc: us,
	}
}

func (s *service) GetMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := s.uc.GetMovies()

	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	respondWithJSON(w, http.StatusOK, movies)
}

func (s *service) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movie types.Movie

	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload, Movie name and Movie Year is required"})
		return
	}

	if movie.MovieName == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"message": "Movie name is required"})
		return
	}
	if movie.MovieYear < YearOfTheFirstEverMovie {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"message": "Movie year is required and must be more than 1895"})
	} else {

		fmt.Println("Inserting movie into DB")

		id, err := s.uc.CreateMovie(movie)

		if err != nil {
			respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		respondWithJSON(w, http.StatusCreated, map[string]string{"status": "created", "id": id.String()})
	}
}

func (s *service) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := uuid.Parse(params["id"])
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err = s.uc.DeleteMovie(id)
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Item with id %v not found", id)})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "deleted", "id": id.String()})

}

func (s *service) GetMovieById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := uuid.Parse(params["id"])
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	movie, err := s.uc.GetMovieById(id)
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Item with id %v not found", id)})
		return
	}

	respondWithJSON(w, http.StatusOK, movie)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
