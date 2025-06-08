package main

import (
	"context"
	"net/http"

	"github.com/dalemusser/gowebcore/asset"
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
	r.Mount("/assets", asset.Handler())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{
			"CSS": asset.Path("app.css"),
			"JS":  asset.Path("alpine.js"),
		}
		render.Render(w, r, "home.html", data)
	})

	srv := server.New(c.Base, r)
	_ = server.Graceful(context.Background(), srv)
}
