package server

import (
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/jialeicui/feedpilot/pkg/config"
	"github.com/jialeicui/feedpilot/pkg/platform"
)

type Server struct {
	cfg    *config.Config
	engine *gin.Engine

	mu        sync.RWMutex
	platforms map[string]platform.Platform
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg:       cfg,
		platforms: make(map[string]platform.Platform),
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

func (s *Server) RegisterPlatform(p platform.Platform) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.platforms[p.Name()]; ok {
		return
	}
	s.platforms[p.Name()] = p
}
