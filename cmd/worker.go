package cmd

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/dalemusser/gowebcore/logger"
	"github.com/dalemusser/gowebcore/queue"
	redisq "github.com/dalemusser/gowebcore/queue/redis"
	"github.com/dalemusser/gowebcore/tasks"
	"github.com/spf13/cobra"
)

// workerCmd starts background queue consumers and periodic tasks.
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start background queue consumers / cron tasks",
	RunE: func(cmd *cobra.Command, _ []string) error {
		// graceful shutdown on Ctrl-C / SIGTERM
		ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		mgr := tasks.New()

		// ── Queue consumer example (Redis Streams) ──────────────────────────
		if Cfg.Redis.URL != "" {
			rdb, err := redisq.New(Cfg.Redis.URL)
			if err != nil {
				return err
			}
			redisq.NewConsumer(rdb, "email", "workers", emailHandler).Start(mgr)
		}

		// ── Cron-style periodic task: every 5 minutes ───────────────────────
		mgr.Add(func(ctx context.Context) error {
			ticker := time.NewTicker(5 * time.Minute)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-ticker.C:
					if err := purgeCache(ctx); err != nil {
						logger.Instance().Error("purge cache", "err", err)
					}
				}
			}
		})

		// run until ctx cancelled, then wait for tasks to finish
		mgr.Start(ctx)
		mgr.Wait()
		return nil
	},
}

// ----------------------------------------------------------------------
// Example handlers (replace with real logic in your service)
// ----------------------------------------------------------------------

func emailHandler(ctx context.Context, j *queue.Job) error {
	// TODO: process j.Data, send email, etc.
	return nil
}

func purgeCache(ctx context.Context) error {
	// TODO: remove expired items, clear local cache, etc.
	return nil
}
