package cache

import (
	"context"

	"github.com/bibi-ic/mata/internal/models"
)

type Cache interface {
	Set(ctx context.Context, key string, value *models.Meta) error
	Get(ctx context.Context, key string) (*models.Meta, error)
}
