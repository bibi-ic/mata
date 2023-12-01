package controller

import (
	"net/http"
	gourl "net/url"

	"github.com/bibi-ic/mata/api"
	"github.com/bibi-ic/mata/model"
	"github.com/gin-gonic/gin"
)

type MataController struct {
	Key string
}

// Retrieve result for MataController with link input
func (m MataController) Retrieve(c *gin.Context) {
	u := c.Query("url")
	_, err := gourl.ParseRequestURI(u)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(ErrInvalidLink))
		return
	}

	r, err := api.NewIframelyRequest(u, m.Key)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res, err := api.IframelyResponse(r)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	meta := model.Meta{}
	err = meta.UnmarshalJSON(res)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, meta)
}
