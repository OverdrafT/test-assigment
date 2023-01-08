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
	router.HandleFunc("/movies/{movieid}", m.DeleteMovie).Methods("DELETE")

	router.HandleFunc("/upload", f.UploadFile)

	return router
}
