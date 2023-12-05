package server

import (
	"github.com/bibi-ic/mata/config"
	db "github.com/bibi-ic/mata/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config config.Config
	store  db.Store
	router *gin.Engine
}

// New Server HTTP
func New(cfg config.Config, store db.Store) *Server {
	s := &Server{
		config: cfg,
		store:  store,
	}

	s.newRouter()
	return s
}

func (s *Server) Start() error {
	address := s.config.Server.Address + ":" + s.config.Server.Port
	return s.router.Run(address)
}
