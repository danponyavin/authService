package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	Server *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(handler http.Handler) error {
	s.Server = &http.Server{
		Addr:           ":8000",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return s.Server.ListenAndServe()
}

func (s *Server) Shutdown() error {
	return s.Server.Shutdown(context.Background())
}
