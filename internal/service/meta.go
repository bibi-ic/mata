package service

import (
	"encoding/json"
	"net/http"

	"github.com/bibi-ic/mata/api"
	"github.com/bibi-ic/mata/internal/cache"
	"github.com/bibi-ic/mata/internal/datastruct"
	db "github.com/bibi-ic/mata/internal/db/sqlc"
	"github.com/bibi-ic/mata/internal/dto"
	"github.com/bibi-ic/mata/internal/status"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type MetaService interface {
	Retrieve(c *gin.Context, u string) (*datastruct.Meta, status.Status)
}

type metaService struct {
	cache cache.MataCache
	store db.Store
}

func NewMetaService(cache cache.MataCache, store db.Store) MetaService {
	return &metaService{
		cache: cache,
		store: store,
	}
}

func (m *metaService) Retrieve(c *gin.Context, u string) (*datastruct.Meta, status.Status) {
	var err error

	key, err := m.store.GetAPITx(c)
	if err != nil {
		return nil, status.Status{
			Code:  http.StatusInternalServerError,
			Error: err,
		}
	}

	mDto := new(dto.Meta)
	metaCached, err := m.cache.Get(c, u)
	switch {
	case err == redis.Nil || metaCached == nil:
		r, err := api.NewIframelyRequest(u, key)
		if err != nil {
			return nil, status.Status{
				Code:  http.StatusInternalServerError,
				Error: err,
			}
		}

		res, err := api.IframelyResponse(r)
		if err != nil {
			return nil, status.Status{
				Code:  http.StatusInternalServerError,
				Error: err,
			}
		}

		err = json.Unmarshal(res, mDto)
		if err != nil {
			return nil, status.Status{
				Code:  http.StatusUnprocessableEntity,
				Error: err,
			}
		}

		meta := new(datastruct.Meta)
		meta.Parse(*mDto)

		err = m.cache.Set(c, mDto.URL, meta)
		if err != nil {
			return nil, status.Status{
				Code:  http.StatusInternalServerError,
				Error: err,
			}
		}

		return meta, status.Status{
			Code:  http.StatusOK,
			Error: nil,
		}

	case err != nil:
		return nil, status.Status{
			Code:  http.StatusInternalServerError,
			Error: err,
		}
	}

	return metaCached, status.Status{
		Code:  http.StatusOK,
		Error: nil,
	}
}
