package cmd

import (
	"context"

	"github.com/dalemusser/gowebcore/server"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := chi.NewRouter()
		srv := server.New(Cfg, r)
		return server.Serve(context.Background(), srv, Cfg.CertFile, Cfg.KeyFile)
	},
}
