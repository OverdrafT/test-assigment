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
		zap.S().Fatal("failed to connect to DB")
	}
	zap.S().Info("connected to db successfully")

	err = db.AutoMigrate(&types.Author{})
	if err != nil {
		zap.S().Fatal("failed to create Authors Table")
	}

	err = db.AutoMigrate(&types.Movie{})
	if err != nil {
		zap.S().Fatal("failed to create Movies table")
	}
	zap.S().Info("migration complete")

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

func (s *service) GetMovies() (movies []types.Movie, err error) {
	result := s.db.Find(&movies)
	if result.Error != nil {
		zap.S().Error("no such items", result.Error)
		return nil, result.Error
	}

	return movies, nil
}

func (s *service) CreateMovie(movie types.Movie) (uuid.UUID, error) {
	result := s.db.Create(&movie)
	if result.Error != nil {
		zap.S().Error("failed to create movie", result.Error)
		return uuid.Nil, result.Error
	}

	return movie.ID, nil
}

func (s *service) DeleteMovie(ID uuid.UUID) error {
	result := s.db.Delete(&types.Movie{ID: ID})
	if result.Error != nil {
		zap.S().Error("failed to delete movie", result.Error)
		return result.Error
	}
	return nil
}

func (s *service) GetMovieById(ID uuid.UUID) (types.Movie, error) {
	movie := types.Movie{ID: ID}
	result := s.db.First(&movie)
	if result.Error != nil {
		zap.S().Error("failed to get movie", result.Error)
		return movie, result.Error
	}
	return movie, nil

}

func (s *service) GetMovieByAuthor(AuthorID uuid.UUID) (movies []types.Movie, err error) {
	err = s.db.Find(&movies, "author_id = ?", AuthorID).Error
	if err != nil {
		zap.S().Error("failed to get movies by author", err)
		return movies, err
	}

	return movies, nil
}

func (s *service) GetAuthors() (authors []types.Author, err error) {
	result := s.db.Find(&authors)
	if result.Error != nil {
		zap.S().Error("failed to get authors", result.Error)
		return nil, result.Error
	}

	return authors, nil
}

func (s *service) CreateAuthor(author types.Author) (uuid.UUID, error) {
	result := s.db.Create(&author)
	if result.Error != nil {
		zap.S().Error("failed to create author", result.Error)
		return uuid.Nil, result.Error
	}

	return author.AuthorID, nil
}

func (s *service) DeleteAuthor(ID uuid.UUID) error {
	result := s.db.Delete(&types.Author{AuthorID: ID})
	if result.Error != nil {
		zap.S().Error("failed to delete author", result.Error)
		return result.Error
	}
	return nil
}

func (s *service) GetAuthorById(ID uuid.UUID) (types.Author, error) {
	author := types.Author{AuthorID: ID}
	result := s.db.First(&author)
	if result.Error != nil {
		zap.S().Error("failed to get author", result.Error)
		return author, result.Error
	}
	return author, nil

}

func (s *service) GetAuthorWithMovies(AuthorID uuid.UUID) (author []types.Author, err error) {
	err = s.db.Preload("Movie").Find(&author, "author_id = ?", AuthorID).Error
	if err != nil {
		zap.S().Error("failed to get author with movies", err)
		return nil, err
	}

	return author, err
}

func (s *service) TriggerPanic() {
	arr := []int{1}
	fmt.Println(arr[100]) //nil pointer exception
}
