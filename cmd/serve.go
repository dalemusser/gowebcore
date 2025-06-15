// cmd/serve.go
package cmd

import (
	"context"
	"net/http"

	"github.com/dalemusser/gowebcore/server"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
)

// serveCmd starts the HTTP server.
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	RunE: func(cmd *cobra.Command, _ []string) error {
		// build router
		r := chi.NewRouter()
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("gowebcore up\n"))
		})

		// NOTE: pass Cfg.Base (not Cfg) to satisfy server.New
		srv := server.New(Cfg.Base, r)

		// run until context cancelled
		return server.Serve(context.Background(), srv, Cfg.CertFile, Cfg.KeyFile)
	},
}
