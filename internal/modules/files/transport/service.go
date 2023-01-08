package transport

import (
	"fmt"
	"io/ioutil"
	"net/http"

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

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")

	if err != nil {
		zap.S().Error(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	tempFile, err := ioutil.TempFile("upload-files", "upload-*.txt")

	if err != nil {
		zap.S().Error(err)
		return
	}

	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		zap.S().Error(err)
		return
	}

	tempFile.Write(fileBytes)

	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
