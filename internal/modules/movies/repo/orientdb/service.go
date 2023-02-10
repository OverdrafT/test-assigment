package orientdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"test-assigment/internal/config"
	"test-assigment/internal/modules/movies/types"
	"time"

	"github.com/google/uuid"

	"go.uber.org/zap"
)

type service struct {
	client       *http.Client
	orientDBHost string
	orientDBName string
}

func New(cfg *config.Config) *service {
	s := service{
		client:       new(http.Client),
		orientDBHost: cfg.ORIENT_DB_HOST,
		orientDBName: cfg.ORIENT_DB_NAME,
	}

	time.Sleep(12 * time.Second)
	err := s.PingDB()
	if err != nil {
		zap.S().Fatal("Connection to db failed, db is down:", err)
	}

	err = s.CheckDB()
	if err != nil {
		zap.S().Error(err)
		err = s.SetupDB()
		if err != nil {
			zap.S().Fatal("failed to setup db", err)
		}
	}

	return &s
}

func (s *service) TriggerPanic() {
	arr := []int{1}
	fmt.Println(arr[100]) //nil pointer exception
}

func (s *service) PingDB() error {
	body, err := s.SendRequest("GET", fmt.Sprintf("http://%s:2480/listDatabases", s.orientDBHost), nil)
	if err != nil {
		zap.S().Error(err)
		return err
	}
	fmt.Println(bytes.NewBuffer(body).String())

	return nil
}

func (s *service) CheckDB() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:2480/connect/%s", s.orientDBHost, s.orientDBName), nil)
	if err != nil {
		zap.S().Info("Incorrect request", err)
		return err
	}

	req.SetBasicAuth("root", "root")
	req.Header.Set("Accept-Encoding", "gzip,deflate")
	req.Header.Set("Content-Length", "52098")

	resp, err := s.client.Do(req)
	if err != nil {
		zap.S().Error("failed to send request:", err)
		return err
	}

	if resp.StatusCode == 401 {
		return errors.New("no such db")
	}

	return nil
}

func (s *service) SetupDB() error {
	_, err := s.SendRequest("POST", fmt.Sprintf("http://%s:2480/database/%s/plocal", s.orientDBHost, s.orientDBName), nil)
	if err != nil {
		zap.S().Error("failed to create db,", err)
		fmt.Println(err)
		return err
	}
	zap.S().Info("DB created")

	body, err := s.SendRequest("GET", fmt.Sprintf("http://%s:2480/listDatabases", s.orientDBHost), nil)

	fmt.Println(bytes.NewBuffer(body).String())

	//create vertex Movie
	url := fmt.Sprintf("http://%s:2480/command/%s/sql/create class Movies extends V",
		s.orientDBHost, s.orientDBName,
	)

	_, err = s.SendRequest("POST", url, nil)

	if err != nil {
		zap.S().Error("failed to create vertex Movie", err)
		return err
	}

	//create vertex Author
	url = fmt.Sprintf("http://%s:2480/command/%s/sql/create class Author extends V",
		s.orientDBHost, s.orientDBName,
	)

	_, err = s.SendRequest("POST", url, nil)

	if err != nil {
		zap.S().Error("failed to create vertex Author", err)
		return err
	}

	//create edge IsAuthor
	url = fmt.Sprintf("http://%s:2480/command/%s/sql/create class IsAuthor extends E",
		s.orientDBHost, s.orientDBName,
	)

	_, err = s.SendRequest("POST", url, nil)

	if err != nil {
		zap.S().Error("failed to create edge IsAuthor", err)
		return err
	}

	return nil
}

func (s *service) SendRequest(method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		zap.S().Info("Incorrect removiequest", err)
		return nil, err
	}

	req.SetBasicAuth("root", "root")
	req.Header.Set("Accept-Encoding", "gzip,deflate")
	req.Header.Set("Content-Length", "52098")

	resp, err := s.client.Do(req)
	if err != nil {
		zap.S().Error("failed to send request:", err)
		return nil, err
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.S().Error("failed to read response body: ", err)
		return nil, err
	}

	return contents, nil
}

func (s *service) GetMovies() ([]types.Movie, error) {

	body, err := s.SendRequest("GET",
		fmt.Sprintf("http://%s:2480/query/%s/sql/select from Movies", s.orientDBHost, s.orientDBName), nil)

	var data orientdbGetMoviesResponse

	err = json.Unmarshal(body, &data)
	if err != nil {
		zap.S().Error("failed to unmarshal response body:", err)
		return nil, err
	}

	return data.Result, nil
}

func (s *service) CreateAuthor(author types.Author) (uuid.UUID, error) {

	create := orientdbCreateAuthorResponse{
		Author: author,
		Class:  "Author",
	}

	create.Author.AuthorID = uuid.New()

	url := fmt.Sprintf(
		"http://%s:2480/command/%s/sql/INSERT INTO Author SET ID = '%s', FirstName = '%s', LastName = '%s'",
		s.orientDBHost, s.orientDBName, create.AuthorID, author.FirstName, author.LastName,
	)

	_, err := s.SendRequest("POST", url, nil)
	if err != nil {
		zap.S().Error("failed to post to db:", err)
		return uuid.Nil, err
	}

	return create.AuthorID, err
}

func (s *service) DeleteMovie(ID uuid.UUID) error {
	_ = types.Movie{ID: ID}

	url := fmt.Sprintf("http://%s:2480/command/%s/sql/delete vertex Movies where ID='%s'", s.orientDBHost, s.orientDBName, ID)

	_, err := s.SendRequest("POST", url, nil)
	if err != nil {

		zap.S().Error("failed to send request", err)
		return err
	}

	return nil
}

func (s *service) GetMovieById(ID uuid.UUID) (types.Movie, error) {

	movie := types.Movie{ID: ID}

	url := fmt.Sprintf("http://%s:2480/query/%s/sql/select from Movies where ID = '%s'", s.orientDBHost, s.orientDBName, ID)

	body, err := s.SendRequest("GET", url, nil)
	if err != nil {
		zap.S().Error("failed to send request", err)
		return movie, err
	}

	var data orientdbGetMovieByID

	err = json.Unmarshal(body, &data)
	if err != nil {
		zap.S().Error("failed to unmarshal response body:", err)
		return movie, err
	}

	return data.Result[0], nil
}

func (s *service) GetMovieByAuthor(ID uuid.UUID) ([]types.Movie, error) {

	url := fmt.Sprintf(
		"http://%s:2480/query/%s/sql/select from Movies where in('IsAuthor').ID = '%s'",
		s.orientDBHost, s.orientDBName, ID,
	)

	body, err := s.SendRequest("GET", url, nil)

	var data orientdbGetMoviesResponse

	err = json.Unmarshal(body, &data)
	if err != nil {
		zap.S().Error("failed to unmarshal response body:", err)
		return nil, err
	}

	return data.Result, nil
}

func (s *service) GetAuthors() ([]types.Author, error) {

	body, err := s.SendRequest("GET",
		fmt.Sprintf("http://%s:2480/query/%s/sql/select from Author", s.orientDBHost, s.orientDBName), nil)

	var data orientdbGetAuthorsResponse

	err = json.Unmarshal(body, &data)
	if err != nil {
		zap.S().Error("failed to unmarshal response body:", err)
		return nil, err
	}

	return data.Result, nil
}

func (s *service) CreateMovie(movie types.Movie) (uuid.UUID, error) {

	create := orientdbCreateMovieRequest{
		Movie: movie,
		Class: "Movies",
	}

	url := fmt.Sprintf("http://%s:2480/query/%s/sql/select from Author where ID = '%s'", s.orientDBHost, s.orientDBName, movie.AuthorID)
	body, err := s.SendRequest("GET", url, nil)
	if err != nil {
		zap.S().Error("failed to send request")
		return uuid.Nil, err
	}

	var data orientdbGetAuthorsResponse

	err = json.Unmarshal(body, &data)
	if err != nil {
		zap.S().Error("failed to unmarshal response body:", err)
		return uuid.Nil, err
	}

	if len(data.Result) == 0 {
		zap.S().Error("Invalid Author ID")
		return uuid.Nil, errors.New("invalid Author ID")
	}

	create.Movie.ID = uuid.New()

	url = fmt.Sprintf(
		"http://%s:2480/command/%s/sql/INSERT INTO Movies SET ID = '%s', MovieYear = %d, MovieName = '%s', AuthorID = '%s'",
		s.orientDBHost, s.orientDBName, create.ID, movie.MovieYear, movie.MovieName, movie.AuthorID,
	)

	_, err = s.SendRequest("POST", url, nil)
	if err != nil {
		zap.S().Error("failed to post to db:", err)
		return uuid.Nil, err
	}

	url = fmt.Sprintf(
		"http://%s:2480/command/%s/sql/create edge IsAuthor from (select from Author where ID = '%s') to (select from Movies where ID='%s')",
		s.orientDBHost, s.orientDBName, movie.AuthorID, create.ID,
	)
	_, err = s.SendRequest("POST", url, nil)
	if err != nil {
		zap.S().Error("failed to post to db:", err)
		return uuid.Nil, err
	}

	return create.ID, err
}

func (s *service) DeleteAuthor(ID uuid.UUID) error {
	_ = types.Author{AuthorID: ID}

	url := fmt.Sprintf("http://%s:2480/command/%s/sql/delete vertex Author where ID='%s'", s.orientDBHost, s.orientDBName, ID)

	_, err := s.SendRequest("POST", url, nil)
	if err != nil {

		zap.S().Error("failed to send request", err)
		return err
	}

	return nil
}

func (s *service) GetAuthorWithMovies(AuthorID uuid.UUID) (author []types.Author, err error) {
	url := fmt.Sprintf(
		"http://%s:2480/query/%s/sql/select from Movies where in('IsAuthor').ID = '%s'",
		s.orientDBHost, s.orientDBName, AuthorID,
	)

	body, err := s.SendRequest("GET", url, nil)

	var moviedata orientdbGetMoviesResponse

	err = json.Unmarshal(body, &moviedata)
	if err != nil {
		zap.S().Error("failed to unmarshal response body:", err)
		return author, err
	}

	url = fmt.Sprintf(
		"http://%s:2480/query/%s/sql/select from Author where ID = '%s'",
		s.orientDBHost, s.orientDBName, AuthorID)

	body, err = s.SendRequest("GET", url, nil)
	if err != nil {
		zap.S().Error("failed to send request", err)
		return author, err
	}

	var data orientdbGetAuthorsResponse

	err = json.Unmarshal(body, &data)
	if err != nil {
		zap.S().Error("failed to unmarshal response body:", err)
		return author, err
	}

	data.Result[0].AuthorID = AuthorID
	data.Result[0].Movie = moviedata.Result

	return data.Result, nil
}
