package postgres

import (
	"fmt"
	"test-assigment/internal/config"
	"test-assigment/internal/modules/movies/types"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupDB(cfg *config.Config) *gorm.DB {

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s  sslmode=disable", cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_NAME, cfg.DB_HOST)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		zap.S().Fatal("Failed to open DB")
	}
	zap.S().Info("Database conncted")

	err = db.AutoMigrate(&types.Movie{})
	if err != nil {
		zap.S().Fatal("Can`t create table")
	}

	return db
}

type service struct {
	db *gorm.DB
}

func New(cfg *config.Config) *service {
	return &service{
		db: setupDB(cfg),
	}
}

// func (s *service) CloseDB() {
// 	s.db.Close()
// }

func (s *service) GetMovies() (movies *gorm.DB, err error) {
	return s.db.Find(&movies), nil
}

func (s *service) CreateMovie(movie types.Movie) (DBid string, err error) {
	panic("implement me!")
}

func (s *service) DeleteMovie(movieID string) error {
	panic("implement me!")
}
