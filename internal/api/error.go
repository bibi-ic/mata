package api

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidLink = errors.New("link is invalid")
)

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
