package server

import (
	"github.com/bibi-ic/mata/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config config.Config
	router *gin.Engine
}

// New Server HTTP
func New() (*Server, error) {
	c, err := config.Load()
	if err != nil {
		return nil, err
	}

	s := &Server{
		config: c,
	}

	s.newRouter()
	return s, err
}

func (s *Server) Start() error {
	address := s.config.Server.Address + ":" + s.config.Server.Port
	return s.router.Run(address)
}
