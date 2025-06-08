package tasks

import (
	"context"
	"testing"
	"time"
)

func TestManager(t *testing.T) {
	m := New()
	m.Add(func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	})
	ctx, cancel := context.WithCancel(context.Background())
	m.Start(ctx)
	cancel()
	done := make(chan struct{})
	go func() { m.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("tasks did not shut down")
	}
}
