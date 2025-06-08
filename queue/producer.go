package queue

import "context"

// Producer publishes jobs to a stream.
type Producer interface {
	Publish(ctx context.Context, stream string, data []byte) error
}
