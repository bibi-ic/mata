package server

import (
	"net/http"

	"github.com/bibi-ic/mata/cache"
	"github.com/bibi-ic/mata/config"
	db "github.com/bibi-ic/mata/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config config.Config
	store  db.Store
	cache  cache.MataCache
	router *gin.Engine
}

// New Server HTTP
func New(cfg config.Config, store db.Store, cache cache.MataCache) *Server {
	s := &Server{
		config: cfg,
		store:  store,
		cache:  cache,
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
