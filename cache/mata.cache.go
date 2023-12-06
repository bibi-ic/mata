package cache

import (
	"context"

	"github.com/bibi-ic/mata/model"
)

type MataCache interface {
	Set(ctx context.Context, key string, value *model.Meta) error
	Get(ctx context.Context, key string) (*model.Meta, error)
}
