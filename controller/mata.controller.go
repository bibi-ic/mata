package controller

import (
	"net/http"
	gourl "net/url"

	"github.com/bibi-ic/mata/api"
	"github.com/bibi-ic/mata/cache"
	db "github.com/bibi-ic/mata/db/sqlc"
	"github.com/bibi-ic/mata/model"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type MataController struct {
	Store db.Store
	Cache cache.MataCache
	Key   string
}

// Retrieve result for MataController with link input
func (m MataController) Retrieve(c *gin.Context) {
	u := c.Query("url")
	_, err := gourl.ParseRequestURI(u)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(ErrInvalidLink))
		return
	}

	m.Key, err = m.Store.GetAPITx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	meta := new(model.Meta)
	d, err := m.Cache.Get(c, u)
	switch {
	case err == redis.Nil || d == nil:
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

		err = meta.ParseJSON(res)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
			return
		}

		err = m.Cache.Set(c, meta.URL, meta)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		c.JSON(http.StatusOK, *meta)
		return

	case err != nil:
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, d)
}
