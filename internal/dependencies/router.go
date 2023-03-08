package dependencies

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	files "test-assigment/internal/modules/files/transport"
	movies "test-assigment/internal/modules/movies/transport"
)

func InitRouter(m movies.Transport, f files.Transport) *mux.Router {
	router := mux.NewRouter()
	router.Use(LoggingMiddleware)

	router.HandleFunc("/movies/", m.GetMovies).Methods("GET")
	router.HandleFunc("/movies/", m.CreateMovie).Methods("POST")
	router.HandleFunc("/movies/{id}", m.DeleteMovie).Methods("DELETE")
	router.HandleFunc("/movies/{id}", m.GetMovieById).Methods("GET")

	router.HandleFunc("/movies/authors/{id}", m.GetMovieByAuthor).Methods("GET")

	router.HandleFunc("/authors/", m.GetAuthors).Methods("GET")
	router.HandleFunc("/authors/", m.CreateAuthor).Methods("POST")
	router.HandleFunc("/authors/{id}", m.DeleteAuthor).Methods("DELETE")
	router.HandleFunc("/authors/{id}", m.GetAuthorById).Methods("GET")
	router.HandleFunc("/authors/movies/{id}", m.GetAuthorWithMovies).Methods("GET")

	router.HandleFunc("/triggerpanic/", m.TriggerPanic).Methods("GET")

	router.HandleFunc("/upload", f.UploadFile)

	return router
}

func FileLogger(filename string) *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	logFile, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		t := time.Now()

		switch r.Method {
		case http.MethodGet:
			WriteLogger(r, t)
		case http.MethodPost:
			WriteLogger(r, t)
		case http.MethodDelete:
			WriteLogger(r, t)
		case http.MethodPut:
			WriteLogger(r, t)
		case http.MethodPatch:
			WriteLogger(r, t)
		}

		next.ServeHTTP(w, r)
	})
}

func WriteLogger(r *http.Request, t time.Time) {
	filename := fmt.Sprintf("./logs/%s/%d.%s.%d-%d.log", r.Method, t.Year(), t.Month(), t.Day(), t.Hour())
	logger := FileLogger(filename)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	logger.Info(r.Method, zap.String("Data:", string(body)), zap.String("UserAgent:", r.UserAgent()), zap.Strings("Header:", r.Header["User-Agent"]))

}
