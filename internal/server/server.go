package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(address string, handler http.Handler) *Server {
	return &Server{httpServer: &http.Server{
		Addr:           address,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	},
	}
}

func (s *Server) Run() {
	if err := s.httpServer.ListenAndServe(); err != nil {
		log.Fatalf("error occured while running server: %s", err.Error())
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
