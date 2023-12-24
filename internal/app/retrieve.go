package app

import (
	"net/http"

	gourl "net/url"

	"github.com/gin-gonic/gin"
)

func (m *Controller) Retrieve(c *gin.Context) {
	u := c.Query("url")
	_, err := gourl.ParseRequestURI(u)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(ErrInvalidLink))
		return
	}

	meta, status := m.metaController.Retrieve(c, u)
	if status.Error != nil {
		c.AbortWithStatusJSON(status.Code, errorResponse(status.Error))
		return
	}

	c.JSON(status.Code, meta)
}
