package natsq

import (
	"context"

	"github.com/nats-io/nats.go"
)

type Producer struct {
	js nats.JetStreamContext
}

func NewProducer(js nats.JetStreamContext) *Producer { return &Producer{js} }

func (p *Producer) Publish(ctx context.Context, subject string, data []byte) error {
	_, err := p.js.PublishMsg(&nats.Msg{
		Subject: subject,
		Data:    data,
		Header:  nats.Header{"Ctx": []string{ctx.Value("request_id").(string)}},
	})
	return err
}
