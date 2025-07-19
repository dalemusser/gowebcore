package main

import (
	"context"
	"log"
	"time"

	"github.com/dalemusser/gowebcore/config"
	"github.com/dalemusser/gowebcore/logger"
	"github.com/dalemusser/gowebcore/server"
	"github.com/dalemusser/gowebcore/tasks"

	"github.com/go-chi/chi/v5"
)

type cfg struct{ config.Base }

func main() {
	var c cfg
	_ = config.Load(&c)
	logger.Init(c.LogLevel)

	// --- set up tasks --------------------
	mgr := tasks.New()
	mgr.Add(tasks.Wrap("ticker", tickerTask))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mgr.Start(ctx)

	// --- http server ---------------------
	r := chi.NewRouter()
	srv := server.New(c.Base, r)

	go func() {
		if err := server.Graceful(ctx, srv); err != nil {
			log.Println(err)
		}
		cancel() // stop tasks when server exits
	}()

	mgr.Wait() // block until tasks finish
}

func tickerTask(ctx context.Context) error {
	t := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			logger.Info("tick")
		}
	}
}
