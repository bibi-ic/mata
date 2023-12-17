package cache

import (
	"context"

	"github.com/bibi-ic/mata/internal/datastruct"
)

type MataCache interface {
	Set(ctx context.Context, key string, value *datastruct.Meta) error
	Get(ctx context.Context, key string) (*datastruct.Meta, error)
}
