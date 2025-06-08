package redis

import (
	"context"
	"time"

	"github.com/dalemusser/gowebcore/logger"
	"github.com/dalemusser/gowebcore/queue"
	"github.com/dalemusser/gowebcore/tasks"
	"github.com/redis/go-redis/v9"
)

type Consumer struct {
	rdb      *redis.Client
	stream   string
	group    string
	handler  queue.HandlerFunc
	maxRetry int
}

// NewConsumer creates/joins a consumer-group on stream.
func NewConsumer(rdb *redis.Client, stream, group string, h queue.HandlerFunc) *Consumer {
	_ = rdb.XGroupCreateMkStream(context.Background(), stream, group, "$").Err()
	return &Consumer{rdb: rdb, stream: stream, group: group, handler: h, maxRetry: 5}
}

// Start registers a background worker in tasks.Manager.
func (c *Consumer) Start(m *tasks.Manager) {
	m.Add(tasks.Wrap("redis-"+c.stream, c.run))
}

func (c *Consumer) run(ctx context.Context) error {
	for {
		// block for up to 5s waiting for messages
		res, err := c.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    c.group,
			Consumer: "gowebcore",
			Streams:  []string{c.stream, ">"},
			Block:    5 * time.Second,
			Count:    10,
		}).Result()
		if err == redis.Nil || (err != nil && ctx.Err() == nil) {
			continue // timeout; loop again
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		for _, s := range res {
			for _, msg := range s.Messages {
				data := msg.Values["data"].(string)
				job := queue.Job{
					ID:     msg.ID,
					Stream: s.Stream,
					Data:   []byte(data),
				}
				if err := c.handler(ctx, &job); err != nil {
					logger.Instance().Error("job error", "err", err, "id", job.ID)
					// leave pending for retry
					continue
				}
				_ = c.rdb.XAck(ctx, s.Stream, c.group, msg.ID).Err()
				_ = c.rdb.XDel(ctx, s.Stream, msg.ID).Err()
			}
		}
	}
}
