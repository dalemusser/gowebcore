package main

import (
	"context"
        "net/http"

	"github.com/dalemusser/gowebcore/config"
	"github.com/dalemusser/gowebcore/logger"
	"github.com/dalemusser/gowebcore/server"

	"github.com/go-chi/chi/v5"
)

type cfg struct{ config.Base }

func main() {
	var c cfg
	_ = config.Load(&c, config.WithEnvPrefix("EX")) // ignore err for brevity

	logger.Init(c.LogLevel)

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from gowebcore\n"))
	})

	srv := server.New(c.Base, r)
	if err := server.Graceful(context.Background(), srv); err != nil {
		logger.Instance().Error("server error", "err", err)
	}
}
