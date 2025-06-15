package tasks

import (
	"context"
	"testing"
	"time"
)

func TestCronRuns(t *testing.T) {
	m := New()
	count := 0
	m.Cron("@every 100ms", func(ctx context.Context) error {
		count++
		return nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()
	m.Start(ctx)
	m.Wait()
	if count == 0 {
		t.Fatalf("cron did not fire")
	}
}
