package redis

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/dalemusser/gowebcore/queue"
	"github.com/dalemusser/gowebcore/tasks"
)

func TestPublishAndConsume(t *testing.T) {
	redisURL := os.Getenv("REDIS_URL") // e.g. redis://localhost:6379/0
	if redisURL == "" {
		t.Skip("set REDIS_URL to run Redis Streams integration test")
	}

	rdb, err := New(redisURL)
	if err != nil {
		t.Fatalf("redis connect: %v", err)
	}
	defer rdb.Close()

	// clean slate
	_ = rdb.FlushDB(context.Background()).Err()

	prod := NewProducer(rdb)
	_ = prod.Publish(context.Background(), "demo", []byte("hello"))

	done := make(chan struct{})
	handler := func(ctx context.Context, j *queue.Job) error {
		if string(j.Data) != "hello" {
			t.Fatalf("unexpected payload")
		}
		close(done)
		return nil
	}

	mgr := tasks.New()
	NewConsumer(rdb, "demo", "workers", handler).Start(mgr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mgr.Start(ctx)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("message not consumed (ensure Redis is reachable at REDIS_URL)")
	}
	cancel()
	mgr.Wait()
}
