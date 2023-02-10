package dependencies

import (
	"github.com/gorilla/mux"

	files "test-assigment/internal/modules/files/transport"
	movies "test-assigment/internal/modules/movies/transport"
)

func InitRouter(m movies.Transport, f files.Transport) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/movies/", m.GetMovies).Methods("GET")
	router.HandleFunc("/movies/", m.CreateMovie).Methods("POST")
	router.HandleFunc("/movies/{id}", m.DeleteMovie).Methods("DELETE")
	router.HandleFunc("/movies/{id}", m.GetMovieById).Methods("GET")

	router.HandleFunc("/movies/authors/{id}", m.GetMovieByAuthor).Methods("GET")

	router.HandleFunc("/authors/", m.GetAuthors).Methods("GET")
	router.HandleFunc("/authors/", m.CreateAuthor).Methods("POST")
	router.HandleFunc("/authors/{id}", m.DeleteAuthor).Methods("DELETE")
	router.HandleFunc("/authors/movies/{id}", m.GetAuthorWithMovies).Methods("GET")

	router.HandleFunc("/triggerpanic/", m.TriggerPanic).Methods("GET")

	router.HandleFunc("/upload", f.UploadFile)

	return router
}
