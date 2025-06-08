package queue

import "context"   // ← add this line

// Job is the envelope passed to your handler.
type Job struct {
	ID       string
	Stream   string
	Data     []byte
	Attempts int
}

// HandlerFunc processes a Job; return nil to ACK, non-nil to retry.
type HandlerFunc func(ctx context.Context, j *Job) error
