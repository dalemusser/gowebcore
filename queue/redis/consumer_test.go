package redis

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/dalemusser/gowebcore/queue"
	"github.com/dalemusser/gowebcore/tasks"
)

func TestPublishAndConsume(t *testing.T) {
	mr, _ := miniredis.Run()
	defer mr.Close()

	rdb, _ := New("redis://" + mr.Addr())
	prod := NewProducer(rdb)

	// publish one job
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
	NewConsumer(rdb, "demo", "g1", handler).Start(mgr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mgr.Start(ctx)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("message not consumed")
	}
	cancel()
	mgr.Wait()
}
