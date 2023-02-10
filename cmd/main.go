package main

import (
	"log"
	"net/http"
	"test-assigment/internal/config"
	"test-assigment/internal/dependencies"
	"test-assigment/internal/modules/movies/repo/orientdb"

	fileTransport "test-assigment/internal/modules/files/transport"
	transport "test-assigment/internal/modules/movies/transport"
	"test-assigment/internal/modules/movies/usecase"
	"test-assigment/pkg/logger"

	"go.uber.org/zap"
)

func main() {

	cfg, err := config.LoadConfig("config") //config
	if err != nil {
		log.Fatal("Failed to open env", err)
	}

	logger.Init("") //TODO: add loglevel from envvars

	repo := orientdb.New(cfg)

	//repo := postgres.New(cfg) //postgres gorm
	//
	//// repo := postgres.New(cfg) //postgres raw
	//// defer repo.CloseDB()
	//
	//// repo := memory.New() //memory
	//
	uc := usecase.New(repo)

	t := transport.New(uc)

	f := fileTransport.New(uc)

	router := dependencies.InitRouter(t, f)

	zap.S().Info("Server at 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
