package transport

import "net/http"

type Transport interface {
	UploadFile(w http.ResponseWriter, r *http.Request)
}
