package service

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// type ClassInfo struct {
// 	Slot   int `json:"slot"`
// 	Remain int `json:"remain"`
// }

type Cache struct {
	timestamp time.Time
	client    *redis.Client
	cache     map[string]int
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		timestamp: time.Now(),
		client:    client,
	}
}

type CacheInterface interface {
	LoadCache(ctx context.Context) error
	Get(key string) (int, bool)
}

func (c *Cache) LoadCache(ctx context.Context) error {
	keys, err := c.client.Keys(ctx, "*").Result()
	if err != nil {
		return err
	}
	c.cache = make(map[string]int, len(keys))
	for _, key := range keys {
		val, err := c.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		n, err := strconv.Atoi(val)
		if err != nil {
			continue
		}
		c.cache[key] = n
	}
	c.timestamp = time.Now()
	return nil
}

func (c *Cache) Get(key string) (int, bool) {
	v, ok := c.cache[key]
	return v, ok
}
