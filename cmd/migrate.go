package cmd

import "github.com/spf13/cobra"

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("migrations: no-op (implement me)")
		return nil
	},
}
