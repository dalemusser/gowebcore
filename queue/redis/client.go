package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// New returns a Redis client from URL, e.g. "redis://localhost:6379/0".
func New(url string) (*redis.Client, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return redis.NewClient(opts), nil
}
