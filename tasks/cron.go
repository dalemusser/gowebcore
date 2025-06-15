package tasks

import (
	"context"

	"github.com/robfig/cron/v3"
)

// Cron schedules fn according to a CRON spec and registers it
// inside the Manager so it shuts down gracefully with the service.
func (m *Manager) Cron(spec string, fn func(ctx context.Context) error) {
	c := cron.New()
	// run fn inside a goroutine managed by Manager
	_, _ = c.AddFunc(spec, func() { _ = fn(context.Background()) })
	m.Add(func(ctx context.Context) error {
		c.Start()
		<-ctx.Done()
		ctxStop := c.Stop() // wait for running jobs
		<-ctxStop.Done()
		return ctx.Err()
	})
}
