package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bibi-ic/mata/internal/external"
	"github.com/bibi-ic/mata/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	gourl "net/url"
)

var (
	ErrInvalidLink = errors.New("link is invalid")
)

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (m *Server) Retrieve(c *gin.Context) {
	u := c.Query("url")
	_, err := gourl.ParseRequestURI(u)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(ErrInvalidLink))
		return
	}

	key, err := m.store.GetAPITx(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	meta := new(models.Meta)
	metaCached, err := m.cache.Get(c, u)
	switch {
	case err == redis.Nil || metaCached == nil:
		r, err := external.NewIframelyRequest(u, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		res, err := external.IframelyResponse(r)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = json.Unmarshal(res, meta)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		err = meta.Parse()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = m.cache.Set(c, meta.URL, meta)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"msg": "ok"})
		return

	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}
