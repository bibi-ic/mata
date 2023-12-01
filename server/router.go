package server

import (
	"github.com/bibi-ic/mata/controller"
	"github.com/gin-gonic/gin"
)

func (s *Server) newRouter() {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	mC := new(controller.MataController)
	mC.Key = s.config.Iframely.Key
	router.POST("/meta", mC.Retrieve)

	s.router = router

	return
}
