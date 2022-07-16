package handlers

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"

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

		parts := strings.Split(u, "/")
		fmt.Printf("%q", parts)
		if len(parts) > 3 {
			if parts[2] == "gauge" || parts[2] == "counter" {
				if parts[1] == "update" && (parts[4] == "" || parts[4] == "none") {
					logrus.Printf("Got metric %s without id, original url: %s", parts[2], u)
					w.WriteHeader(http.StatusBadRequest)
				} else {
					err := h.repos.Store(u)
					if err != nil {
						logrus.Printf("error while storing metric from url: %v", u)
					}
					w.WriteHeader(http.StatusOK)
				}
			} else {
				logrus.Printf("Uknown type of metrics %s", parts[2])
				w.WriteHeader(http.StatusNotImplemented)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
