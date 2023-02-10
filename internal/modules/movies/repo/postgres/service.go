package postgres

import (
	"fmt"
	"test-assigment/internal/config"
	"test-assigment/internal/modules/movies/types"

	"github.com/google/uuid"

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
	zap.S().Info("Database connected")

	err = db.AutoMigrate(&types.Author{})
	if err != nil {
		zap.S().Fatal("Failed to create Authors Table")
	}

	err = db.AutoMigrate(&types.Movie{})
	if err != nil {
		zap.S().Fatal("Failed to create Movies table")
	}
	zap.S().Info("Migration complete")

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

func (s *service) GetMovies() (movies []types.Movie, err error) {
	result := s.db.Find(&movies)
	if result.Error != nil {
		zap.S().Fatal("No such items")
	}

	return movies, nil
}

func (s *service) CreateMovie(movie types.Movie) (uuid.UUID, error) {
	result := s.db.Create(&movie)
	if result.Error != nil {
		return uuid.Nil, result.Error
	}

	return movie.ID, nil
}

func (s *service) DeleteMovie(ID uuid.UUID) error {
	result := s.db.Delete(&types.Movie{ID: ID})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *service) GetMovieById(ID uuid.UUID) (types.Movie, error) {
	movie := types.Movie{ID: ID}
	result := s.db.First(&movie)
	if result.Error != nil {
		return movie, result.Error
	}
	return movie, nil

}

func (s *service) GetAuthors() (authors []types.Author, err error) {
	result := s.db.Find(&authors)
	if result.Error != nil {
		zap.S().Fatal("No such items")
	}

	return authors, nil
}

func (s *service) CreateAuthor(author types.Author) (uuid.UUID, error) {
	result := s.db.Create(&author)
	if result.Error != nil {
		return uuid.Nil, result.Error
	}

	return author.AuthorID, nil
}

func (s *service) DeleteAuthor(ID uuid.UUID) error {
	result := s.db.Delete(&types.Author{AuthorID: ID})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *service) GetMovieByAuthor(AuthorID uuid.UUID) (movies []types.Movie, err error) {
	err = s.db.Find(&movies, "author_id = ?", AuthorID).Error
	if err != nil {
		return movies, err
	}

	return movies, nil
}

func (s *service) GetAuthorWithMovies(AuthorID uuid.UUID) (author []types.Author, err error) {
	err = s.db.Preload("Movie").Find(&author, "author_id = ?", AuthorID).Error
	if err != nil {
		zap.S().Error("No such author", err)
		return nil, err
	}

	fmt.Println(author)

	return author, err
}

func (s *service) TriggerPanic() {
	arr := []int{1}
	fmt.Println(arr[100]) //nil pointer exception
}
