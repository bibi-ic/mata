package server

import (
	"github.com/bibi-ic/mata/controller"
	"github.com/gin-gonic/gin"
)

func (s *Server) newRouter() {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	mC := &controller.MataController{
		Store: s.store,
		Cache: s.cache,
	}
	router.POST("/meta", mC.Retrieve)

	s.router = router

	return
}
