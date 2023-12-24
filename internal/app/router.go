package app

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) newRouter() {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	s.router = router
}

// RegisterRoutes registers controllers to server
func (s *Server) RegisterRoute(handler *Controller) {
	s.router.POST("/meta", handler.Retrieve)
}
