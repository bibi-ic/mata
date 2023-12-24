package controller

import (
	"encoding/json"
	"net/http"

	"github.com/bibi-ic/mata/api"
	"github.com/bibi-ic/mata/internal/models"
	"github.com/bibi-ic/mata/internal/status"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func (m *metaHandler) Retrieve(c *gin.Context, u string) (*models.Meta, status.Status) {
	var err error

	key, err := m.store.GetAPITx(c)
	if err != nil {
		return nil, status.Status{
			Code:  http.StatusInternalServerError,
			Error: err,
		}
	}

	mDto := new(models.Mata)
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

		meta := new(models.Meta)
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