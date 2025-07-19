package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/dalemusser/gowebcore/logger"
	"github.com/dalemusser/gowebcore/queue"
	"github.com/dalemusser/gowebcore/tasks"
)

type Consumer struct {
	svc        *sqs.Client
	queueURL   string
	dlqURL     string
	handler    queue.HandlerFunc
	maxReceive int
}

func NewConsumer(svc *sqs.Client, queueURL string, handler queue.HandlerFunc) *Consumer {
	return &Consumer{
		svc:        svc,
		queueURL:   queueURL,
		handler:    handler,
		maxReceive: 5,
	}
}

func (c *Consumer) WithDLQ(dlq string) *Consumer { c.dlqURL = dlq; return c }

func (c *Consumer) Start(m *tasks.Manager) {
	m.Add(tasks.Wrap("sqs", c.loop))
}

func (c *Consumer) loop(ctx context.Context) error {
	for {
		out, err := c.svc.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &c.queueURL,
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     20,
		})
		if err != nil || len(out.Messages) == 0 {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			continue
		}
		for _, msg := range out.Messages {
			job := queue.Job{
				ID:   *msg.MessageId,
				Data: []byte(*msg.Body),
			}
			if err := c.handler(ctx, &job); err != nil {
				logger.Error("sqs job error", "err", err, "id", job.ID)
				// let it retry automatically
				continue
			}
			_, _ = c.svc.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      &c.queueURL,
				ReceiptHandle: msg.ReceiptHandle,
			})
		}
	}
}
