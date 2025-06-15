package cmd

import (
	"testing"

	"github.com/spf13/cobra" // ← add this line
)

func TestRootCommands(t *testing.T) {
	for _, c := range []*cobra.Command{serveCmd, migrateCmd, workerCmd} {
		if c.Use == "" {
			t.Fatalf("cmd %#v not initialised", c)
		}
	}
}
