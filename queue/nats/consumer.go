package natsq

import (
	"context"
	"fmt"

	"github.com/dalemusser/gowebcore/logger"
	"github.com/dalemusser/gowebcore/queue"
	"github.com/dalemusser/gowebcore/tasks"
	"github.com/nats-io/nats.go"
)

// Consumer reads jobs from a JetStream durable pull subscription
// and forwards them to the shared queue.HandlerFunc.
type Consumer struct {
	subject  string
	js       nats.JetStreamContext
	handler  queue.HandlerFunc
	durable  string
	pullSize int
}

// NewConsumer creates a pull-based durable consumer.
//   - js       : JetStream context returned by nc.JetStream()
//   - subj     : stream subject, e.g. "email"
//   - durable  : durable consumer name, e.g. "workers"
//   - h        : your job handler
func NewConsumer(js nats.JetStreamContext, subj, durable string, h queue.HandlerFunc) *Consumer {
	return &Consumer{
		subject:  subj,
		js:       js,
		handler:  h,
		durable:  durable,
		pullSize: 10,
	}
}

// Start registers the consumer's run-loop with the tasks manager so it
// shuts down gracefully together with the rest of the service.
func (c *Consumer) Start(m *tasks.Manager) { m.Add(tasks.Wrap("nats-"+c.subject, c.loop)) }

func (c *Consumer) loop(ctx context.Context) error {
	sub, err := c.js.PullSubscribe(c.subject, c.durable)
	if err != nil {
		return err
	}

	for {
		msgs, err := sub.Fetch(c.pullSize, nats.Context(ctx))
		if err != nil && err != nats.ErrTimeout {
			// network hiccup or stream error – log and continue
			logger.Instance().Error("nats fetch", "err", err)
			continue
		}
		if ctx.Err() != nil {
			return ctx.Err() // graceful shutdown
		}

		for _, msg := range msgs {
			meta, _ := msg.Metadata()
			job := queue.Job{
				ID:   fmt.Sprintf("%d", meta.Sequence.Consumer),
				Data: msg.Data,
			}

			if err := c.handler(ctx, &job); err != nil {
				_ = msg.Nak() // retry later
				logger.Instance().Warn("nats job failed", "err", err, "id", job.ID)
			} else {
				_ = msg.Ack() // success
			}
		}
	}
}
