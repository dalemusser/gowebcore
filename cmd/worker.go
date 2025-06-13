package cmd

import "github.com/spf13/cobra"

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start background consumers/tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("worker: no-op (implement me)")
		select {} // placeholder – blocks forever
	},
}
