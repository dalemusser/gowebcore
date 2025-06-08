package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Producer struct{ rdb *redis.Client }

func NewProducer(rdb *redis.Client) *Producer { return &Producer{rdb} }

func (p *Producer) Publish(ctx context.Context, stream string, data []byte) error {
	_, err := p.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		ID:     "*",
		Values: map[string]any{"data": data},
	}).Result()
	return err
}
