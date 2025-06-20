// tasks/cron_test.go
package tasks

import (
	"context"
	"testing"
	"time"
)

func TestCronRuns(t *testing.T) {
	mgr := New()

	fired := make(chan struct{})

	// every 100 ms
	mgr.Cron("@every 100ms", func(ctx context.Context) error {
		select {
		case fired <- struct{}{}: // signal test
		default: // already signalled
		}
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go mgr.Start(ctx)

	select {
	case <-fired:
		// ok
	case <-time.After(2 * time.Second):
		t.Fatal("cron did not fire within 2 s")
	}
}
