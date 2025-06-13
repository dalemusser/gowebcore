package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Producer struct {
	svc      *sqs.Client
	queueURL string
}

func NewProducer(svc *sqs.Client, queueURL string) *Producer {
	return &Producer{svc: svc, queueURL: queueURL}
}

func (p *Producer) Publish(ctx context.Context, _ string, data []byte) error {
	_, err := p.svc.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &p.queueURL,
		MessageBody: aws.String(string(data)),
	})
	return err
}
