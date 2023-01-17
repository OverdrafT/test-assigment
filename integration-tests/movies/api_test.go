package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"test-assigment/internal/config"
	"test-assigment/internal/modules/movies/types"
	"testing"
)

func setupDB(cfg *config.Config) *gorm.DB {

	//dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s  sslmode=disable port=%s",
	//	cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_NAME, cfg.DB_HOST, "5433")

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s  sslmode=disable port=%s",
		"postgres_test", "postgres_test", "postgres_test", "localhost", "5433")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		zap.S().Fatal("Failed to open DB")
	}
	zap.S().Info("Database connected")

	err = db.AutoMigrate(&types.Movie{})
	if err != nil {
		zap.S().Fatal("Can`t create table")
	}
	zap.S().Info("Migration complete")

	return db
}

func TestCreateMovie(t *testing.T) {

	cfg, err := config.LoadConfig("config-test")
	if err != nil {
		log.Fatal("Failed to open env", err)
	}
	db := setupDB(cfg)

	postgres_test, err := db.DB()
	if err != nil {
		zap.S().Fatal(err)
	}
	defer postgres_test.Close()

	movie := types.Movie{
		MovieName: "hello",
		MovieYear: 2000,
	}

	postBody, _ := json.Marshal(movie)
	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post("http://localhost:8081/movies/", "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	fmt.Print(sb)
	db.Delete(types.Movie{ID: movie.ID})
}

func TestGetMovieById(t *testing.T) {

}

func TestGetMovies(t *testing.T) {

	cfg, err := config.LoadConfig("config-test")
	if err != nil {
		log.Fatal("Failed to open env", err)
	}
	db := setupDB(cfg)

	postgres_test, err := db.DB()
	if err != nil {
		zap.S().Fatal(err)
	}
	defer postgres_test.Close()

	movie := types.Movie{
		MovieName: "hello",
		MovieYear: 2000,
	}

	db.Create(movie)

	resp, err := http.Get("http://localhost:8081/movies/")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	log.Printf(sb)
	db.Delete(types.Movie{ID: movie.ID})
}

func TestDeleteMovie(t *testing.T) {

	cfg, err := config.LoadConfig("config-test")
	if err != nil {
		log.Fatal("Failed to open env", err)
	}
	db := setupDB(cfg)

	postgres_test, err := db.DB()
	if err != nil {
		zap.S().Fatal(err)
	}
	defer postgres_test.Close()

	movie := types.Movie{
		MovieName: "hello",
		MovieYear: 2000,
	}

	db.Create(movie)

	url, err := fmt.Printf("localhost:8080/movies/%s", movie.ID)
	if err != nil {
		zap.S().Fatal(err)
	}

	req, err := http.NewRequest("DELETE", strconv.Itoa(url), nil)
	if err != nil {
		zap.S().Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		zap.S().Fatal(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.S().Fatal()
	}

	fmt.Println("response Status :", resp.Status)
	fmt.Println("response Headers :", resp.Header)
	fmt.Println("response Body :", string(respBody))
}
