package server

import (
	"github.com/gin-gonic/gin"

	"github.com/jialeicui/feedpilot/pkg/config"
)

type Server struct {
	cfg *config.Config

	engine *gin.Engine
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Start() error {
	addr := s.cfg.Addr

	s.engine = gin.Default()
	if err := s.registerRouter(); err != nil {
		return err
	}

	return s.engine.Run(addr)
}
