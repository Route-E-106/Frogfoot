package server

import (
	"log/slog"
)

type Server struct {
	logger *slog.Logger
}

func NewServer() *Server {
	return &Server{}

}
