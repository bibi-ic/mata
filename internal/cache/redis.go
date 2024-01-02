package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bibi-ic/mata/internal/models"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
	expire time.Duration
}

func New(client *redis.Client, exp time.Duration) Cache {
	return &redisCache{
		client: client,
		expire: exp,
	}
}

func (cache *redisCache) Set(ctx context.Context, key string, m *models.Meta) error {
	d, _ := json.Marshal(*m)

	err := cache.client.Set(ctx, m.URL, d, time.Duration(m.CacheAge/3600*int64(time.Minute))).Err()
	return err
}

func (cache *redisCache) Get(ctx context.Context, key string) (*models.Meta, error) {
	d, err := cache.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	m := new(models.Meta)
	err = json.Unmarshal([]byte(d), m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
