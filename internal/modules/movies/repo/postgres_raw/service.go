package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"test-assigment/internal/config"
	"test-assigment/internal/modules/movies/types"

	_ "github.com/lib/pq"
)

// const (
// 	DB_USER     = "postgres"
// 	DB_PASSWORD = "postgres"
// 	DB_NAME     = "postgres"
// )

func setupDB(cfg *config.Config) *sql.DB {

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_NAME, cfg.DB_HOST)
	// dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=postgres sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	//dbinfo := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", "postgres", "postgres", "postgres", "postgres")d
	fmt.Println(dbinfo)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal("Failed to run sql.Open", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping db: ", err)
	}

	return db
}

type service struct {
	db *sql.DB
}

func New(cfg *config.Config) *service {
	return &service{
		db: setupDB(cfg),
	}
}

func (s *service) CloseDB() {
	s.db.Close()
}

func (s *service) GetMovies() (movies []types.Movie, err error) {

	fmt.Println("Getting movies...")

	rows, err := s.db.Query("SELECT * FROM movies")

	if err != nil {
		return nil, err
	}
	fmt.Println("Getting movies...2")

	for rows.Next() {
		var id int
		var movieYear int
		var movieName string

		err = rows.Scan(&id, &movieYear, &movieName)

		if err != nil {
			fmt.Println("failed to read movies from db", err)
		}

		fmt.Println("Getting movies...")

		movies = append(movies, types.Movie{MovieYear: movieYear, MovieName: movieName})
	}
	return movies, nil
}

func (s *service) CreateMovie(movie types.Movie) (string, error) {
	var lastInsertID int

	err := s.db.QueryRow("INSERT INTO movies(movieYear, movieName) VALUES($1, $2) returning id;", movie.MovieYear, movie.MovieName).Scan(&lastInsertID)

	if err != nil {
		return "", err
	}

	return strconv.Itoa(lastInsertID), nil
}

func (s *service) DeleteMovie(ID string) error {

	fmt.Println("Deleting movie from DB")

	_, err := s.db.Exec("DELETE FROM movies where ID = $1", ID)

	return err
}
