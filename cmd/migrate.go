package cmd

import (
	"context"
	"time"

	"github.com/dalemusser/gowebcore/db/sql/migrate"
	pg "github.com/dalemusser/gowebcore/db/sql/postgres"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations and exit",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(cmd.Context(), 2*time.Minute)
		defer cancel()

		if Cfg.Postgres.DSN == "" {
			return cmd.Help() // or return an error: "postgres.dsn not set"
		}

		db, err := pg.Open(Cfg.Postgres.DSN)
		if err != nil {
			return err
		}
		defer db.Close()

		return migrate.Run(ctx, db)
	},
}
