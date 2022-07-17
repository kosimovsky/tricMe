package handlers

import (
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

func (h Handler) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		//cType := r.Header.Get("Content-Type")
		//if !compareStr(cType, "text/plain") {
		//	w.WriteHeader(http.StatusUnsupportedMediaType)
		//	_, err := w.Write([]byte("Sorry! Content-Type you sent not accepted"))
		//	if err != nil {
		//		logrus.Printf("wrong Content-Type, couldn't write to response")
		//	}
		//}
		u := r.URL.RequestURI()
		parts := strings.Split(u, "/")
		if len(parts) > 4 {
			if compareStr(parts[2], "gauge") || compareStr(parts[2], "counter") {
				if parts[1] == "update" && (parts[4] == "" || parts[4] == "none") {
					w.WriteHeader(http.StatusBadRequest)
					logrus.Printf("Got metric %s without id, original url: %s", parts[2], u)
				} else {
					err := h.repos.Store(u)
					if err != nil {
						logrus.Printf("error while storing metric from url: %v", u)
					}
					w.WriteHeader(http.StatusOK)
				}
			} else {
				w.WriteHeader(http.StatusNotImplemented)
				logrus.Printf("Uknown type of metrics %s", parts[2])
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func compareStr(str1, str2 string) bool {
	return strings.Compare(str1, str2) == 0
}
