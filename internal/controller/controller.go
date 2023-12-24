package controller

import (
	"github.com/bibi-ic/mata/internal/cache"
	db "github.com/bibi-ic/mata/internal/db/sqlc"
	"github.com/bibi-ic/mata/internal/models"
	"github.com/bibi-ic/mata/internal/status"
	"github.com/gin-gonic/gin"
)

type MetaController interface {
	Retrieve(c *gin.Context, u string) (*models.Meta, status.Status)
}

type metaHandler struct {
	cache cache.MataCache
	store db.Store
}

func NewMetaController(cache cache.MataCache, store db.Store) MetaController {
	return &metaHandler{
		cache: cache,
		store: store,
	}
}
