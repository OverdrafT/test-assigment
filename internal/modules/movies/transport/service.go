package transport

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"test-assigment/internal/modules/movies/types"
	"test-assigment/internal/modules/movies/usecase"
	"time"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"

	"github.com/google/uuid"
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

func FileLogger(filename string) *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	logFile, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger
}

func (s *service) GetMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := s.uc.GetMovies()

	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithJSON(w, r, http.StatusOK, movies)
}

func (s *service) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movie types.Movie

	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{"error": "invalid request payload, Movie name and Movie Year is required"})
		return
	}

	if movie.AuthorID == "" {
		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{"message": "Author ID is required"})
		return
	}
	if len(movie.AuthorID) < 36 {
		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{"message": "Invalid Author ID length"})
		return
	}
	if movie.MovieName == "" {
		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{"message": "Movie name is required"})
		return
	}
	if movie.MovieYear < YearOfTheFirstEverMovie {
		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{"message": "Movie year is required and must be more than 1895"})
		return
	} else {

		fmt.Println("Inserting movie into DB")

		id, err := s.uc.CreateMovie(movie)

		if err != nil {

			respondWithJSON(w, r, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		respondWithJSON(w, r, http.StatusCreated, map[string]string{"status": "created", "id": id.String()})
	}
}

func (s *service) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := validateUuid(w, r)
	if err != nil {
		return
	}

	err = s.uc.DeleteMovie(id)
	if err != nil {
		respondWithJSON(w, r, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Item with id %v not found", id)})
		return
	}

	respondWithJSON(w, r, http.StatusOK, map[string]string{"status": "deleted", "id": id.String()})

}

func (s *service) GetMovieById(w http.ResponseWriter, r *http.Request) {
	id, err := validateUuid(w, r)
	if err != nil {
		return
	}

	movie, err := s.uc.GetMovieById(id)
	if err != nil {
		respondWithJSON(w, r, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Item with id %v not found", id)})
		return
	}

	respondWithJSON(w, r, http.StatusOK, movie)
}

func (s *service) GetMovieByAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := validateUuid(w, r)
	if err != nil {
		return
	}

	movie, err := s.uc.GetMovieByAuthor(id)
	if err != nil {
		respondWithJSON(w, r, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Item with id %v not found", id)})
		return
	}

	respondWithJSON(w, r, http.StatusOK, movie)
}

func (s *service) GetAuthorWithMovies(w http.ResponseWriter, r *http.Request) {
	id, err := validateUuid(w, r)
	if err != nil {
		return
	}

	movie, err := s.uc.GetAuthorWithMovies(id)
	if err != nil {
		respondWithJSON(w, r, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Item with id %v not found", id)})
		return
	}

	respondWithJSON(w, r, http.StatusOK, movie)
}

func (s *service) GetAuthors(w http.ResponseWriter, r *http.Request) {
	authors, err := s.uc.GetAuthors()

	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithJSON(w, r, http.StatusOK, authors)

}

func (s *service) GetAuthorById(w http.ResponseWriter, r *http.Request) {
	id, err := validateUuid(w, r)
	if err != nil {
		return
	}

	movie, err := s.uc.GetAuthorById(id)
	if err != nil {
		respondWithJSON(w, r, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Item with id %v not found", id)})
		return
	}

	respondWithJSON(w, r, http.StatusOK, movie)
}

func (s *service) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var author types.Author

	err := json.NewDecoder(r.Body).Decode(&author)
	if err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{"error": "invalid request payload, First Name and Last Name is required"})
		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if author.FirstName == "" || author.LastName == "" {
		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{"message": "First and Last Name are required"})
		return
	} else {

		fmt.Println("Inserting author into DB")

		id, err := s.uc.CreateAuthor(author)

		if err != nil {
			respondWithJSON(w, r, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		respondWithJSON(w, r, http.StatusCreated, map[string]string{"status": "created", "id": id.String()})
	}
}

func (s *service) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := validateUuid(w, r)
	if err != nil {
		return
	}

	err = s.uc.DeleteAuthor(id)
	if err != nil {
		respondWithJSON(w, r, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Item with id %v not found", id)})
		return
	}

	respondWithJSON(w, r, http.StatusOK, map[string]string{"status": "deleted", "id": id.String()})

}

func (s *service) TriggerPanic(http.ResponseWriter, *http.Request) {
	s.uc.TriggerPanic()
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	if (code / 100) == 2 {
		response, _ := json.Marshal(payload)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_, err := w.Write(response)
		if err != nil {
			zap.S().Error(err)
		}
	} else {
		response, _ := json.Marshal(payload)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_, err := w.Write(response)
		if err != nil {
			zap.S().Error(err)
		}
		t := time.Now()
		filename := fmt.Sprintf("./logs/ERRORS/%d.%s.%d-%d.log", t.Year(), t.Month(), t.Day(), t.Hour())
		logger := FileLogger(filename)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))

		logger.Info(r.Method, zap.String("Error:", string(response)), zap.String("StatusCode:", fmt.Sprintf("%v", code)), zap.String("Data:", string(body)), zap.String("UserAgent:", r.UserAgent()), zap.Strings("Header:", r.Header["User-Agent"]))

	}
}

func validateUuid(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	params := mux.Vars(r)

	if len(params["id"]) < 36 {
		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{"message": "Invalid ID length"})
		return uuid.Nil, errors.New("invalid uuid length")
	}

	id, err := uuid.Parse(params["id"])
	if err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return uuid.Nil, errors.New("failed to parse uuid")
	}

	return id, nil
}
