package transport

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	"test-assigment/internal/modules/movies/usecase"

	"go.uber.org/zap"
)

type service struct {
	f *usecase.Service
}

func New(us *usecase.Service) *service {
	return &service{
		f: us,
	}
}

func (s *service) UploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return
	}

	file, handler, err := r.FormFile("myFile")

	if err != nil {
		zap.S().Error(err)
		return
	}

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	tempFile, err := ioutil.TempFile("", "upload-*.txt")

	if err != nil {
		zap.S().Error(err)
		return
	}

	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {

		}
	}(tempFile)

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		zap.S().Error(err)
		return
	}

	_, err = tempFile.Write(fileBytes)
	if err != nil {
		return
	}

	_, err = fmt.Fprintf(w, "Successfully Uploaded File\n")
	if err != nil {
		return
	}
}
