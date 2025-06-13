package cmd

import (
	"fmt"

	"github.com/dalemusser/gowebcore/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print build information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("commit: %s  built: %s  go: %s\n",
			version.Commit, version.BuildDate, version.GoVersion)
	},
}
