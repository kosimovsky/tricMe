package server

import (
	"context"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(address string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           validate(address),
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
func validate(address string) string {
	return strings.TrimLeft(address, `https://`)
}
