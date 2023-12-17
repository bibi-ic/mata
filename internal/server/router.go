package server

import (
	"github.com/bibi-ic/mata/internal/app"
	"github.com/gin-gonic/gin"
)

func (s *Server) newRouter() {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	s.router = router
}

// RegisterService registers services to server
func (s *Server) RegisterService(service *app.ServiceServer) {
	s.router.POST("/meta", service.Retrieve)
}
