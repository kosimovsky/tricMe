package handlers

import (
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/kosimovsky/tricMe/internal/storage"
)

type Handler struct {
	repos storage.Repositories
}

func NewHandler(repos storage.Repositories) *Handler {
	return &Handler{repos: repos}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if r.Header.Get("Content-Type") != "text/plain" {
			_, err := w.Write([]byte("Sorry! Content-Type you sent not accepted"))
			if err != nil {
				logrus.Printf("wrong Content-Type, couldn't write to response")
			}
		}
		u := r.URL.RequestURI()
		err := h.repos.Store(u)
		if err != nil {
			logrus.Printf("error while storing metric from url: %v", u)
		}
		w.WriteHeader(http.StatusOK)
	}
}
