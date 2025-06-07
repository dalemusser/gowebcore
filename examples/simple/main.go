package main

import (
	"context"
	"net/http"

	"github.com/dalemusser/gowebcore/config"
	"github.com/dalemusser/gowebcore/logger"
	"github.com/dalemusser/gowebcore/render"
	"github.com/dalemusser/gowebcore/server"

	"github.com/go-chi/chi/v5"
)

type cfg struct{ config.Base }

func main() {
	var c cfg
	_ = config.Load(&c, config.WithEnvPrefix("EX"))
	logger.Init(c.LogLevel)

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.Render(w, r, "home.html", nil)
	})
	r.Get("/fragment", func(w http.ResponseWriter, r *http.Request) {
		render.Render(w, r, "fragment.html", map[string]any{
			"Now": "right now",
		})
	})

	srv := server.New(c.Base, r)
	_ = server.Graceful(context.Background(), srv)
}
