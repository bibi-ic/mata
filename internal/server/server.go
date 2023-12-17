package server

import (
	"net/http"

	"github.com/bibi-ic/mata/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config config.Config
	router *gin.Engine
}

func New(cfg config.Config) *Server {
	s := &Server{
		config: cfg,
	}

	s.newRouter()
	return s
}

func (s *Server) Start() error {
	address := s.config.Server.Address + ":" + s.config.Server.Port
	return s.router.Run(address)
}

func (s *Server) DoHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}
