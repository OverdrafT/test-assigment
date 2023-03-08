package main

import (
	"net/http"
	"test-assigment/internal/config"
	"test-assigment/internal/dependencies"
	fileTransport "test-assigment/internal/modules/files/transport"
	"test-assigment/internal/modules/movies/repo/orientdb"
	transport "test-assigment/internal/modules/movies/transport"
	"test-assigment/internal/modules/movies/usecase"
	"test-assigment/pkg/logger"

	"go.uber.org/zap"
)

func main() {

	logger.Init("") //TODO: add loglevel from envvars

	cfg, err := config.LoadConfig("config") //config
	if err != nil {
		zap.S().Fatal("Failed to open env", err)
	}

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
	zap.S().Fatal(http.ListenAndServe(":8080", router))

}
