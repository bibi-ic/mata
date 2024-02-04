package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bibi-ic/mata/internal/external"
	"github.com/bibi-ic/mata/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"

	gourl "net/url"
)

const (
	tracerName = "internal/api"
)

var (
	tracer = otel.Tracer(tracerName)
)

func (m *Server) Retrieve(c *gin.Context) {
	u := c.Query("url")
	_, err := gourl.ParseRequestURI(u)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(ErrInvalidLink))
		return
	}

	ctx, span := tracer.Start(
		c.Request.Context(),
		"retrieve",
		oteltrace.WithAttributes(attribute.String("request.link", u)),
	)
	defer span.End()

	key, err := m.store.GetAPITx(ctx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	meta := new(models.Meta)
	metaCached, err := m.cache.Get(ctx, u)
	switch {
	case err == redis.Nil || metaCached == nil:
		r, err := external.NewIframelyRequest(u, key)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		res, err := external.IframelyResponse(r)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		err = json.Unmarshal(res, meta)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, errorResponse(err))
			return
		}

		err = meta.Parse()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		err = m.cache.Set(c, u, meta)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		span.SetAttributes(attribute.String("metadata.created", fmt.Sprintf("%v", *meta)))
		c.JSON(http.StatusCreated, *meta)
		return

	case err != nil:
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	span.SetAttributes(attribute.String("metadata.found_cache", fmt.Sprintf("%v", *metaCached)))

	c.JSON(http.StatusOK, *metaCached)
}
