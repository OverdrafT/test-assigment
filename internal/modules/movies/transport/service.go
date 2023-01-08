package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"test-assigment/internal/modules/movies/types"
	"test-assigment/internal/modules/movies/usecase"

	"github.com/gorilla/mux"
)

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
	// var response = types.JsonResponse{Type: "success", Data: movies}
	respondWithJSON(w, http.StatusOK, movies)

}

func (s *service) CreateMovie(w http.ResponseWriter, r *http.Request) {
	// movieID := r.FormValue("movieid")
	// movieName := r.FormValue("moviename")
	var movie types.Movie

	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}

	if movie.MovieName == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"message": "Movie name is required"})
		return
	} else {

		fmt.Println("Inserting movie into DB")

		id, err := s.uc.CreateMovie(movie)

		if err != nil {
			respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		respondWithJSON(w, http.StatusCreated, map[string]string{"status": "created", "id": id})
		//response = types.JsonResponse{Type: "success", Message: "The movie has been inserted successfully!"}

		// fmt.Println("Movie ID: " + movieID + "|Movie name: " + movieName + "|DataBase ID: " + id)

	}
}

func (s *service) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]

	if id == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "missing id path param"})
		return
	}

	err := s.uc.DeleteMovie(id)
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Item with id %v not found", id)})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "deleted", "id": id})
	// var response = types.JsonResponse{}

	// if movieID == "" {
	// 	response = types.JsonResponse{Type: "error", Message: "You are missing movieID parameter."}
	// } else {
	// 	s.uc.DeleteMovie(movieID)

	// 	response = types.JsonResponse{Type: "success", Message: "The movie has been deleted successfully!"}
	// }

	// json.NewEncoder(w).Encode(response)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
