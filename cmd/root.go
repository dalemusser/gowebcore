// cmd/root.go
package cmd

import (
	"github.com/dalemusser/gowebcore/config"
	"github.com/dalemusser/gowebcore/logger"
	"github.com/spf13/cobra"
)

// -----------------------------------------------------------------------------
// Configuration model
// -----------------------------------------------------------------------------

// appConfig embeds gowebcore's Base and adds service-specific blocks.
type appConfig struct {
	config.Base

	Redis struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"redis"`

	Postgres struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"postgres"`
}

// -----------------------------------------------------------------------------
// Globals
// -----------------------------------------------------------------------------

var (
	cfgFile string    // --config flag
	Cfg     appConfig // populated in PersistentPreRunE

	rootCmd = &cobra.Command{
		Use:   "gowebsvc",
		Short: "gowebcore service runner (serve | migrate | worker | version)",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// 1. load config (env + file + flags)
			if err := config.Load(&Cfg, config.WithConfigFile(cfgFile)); err != nil {
				return err
			}
			// 2. initialise structured logger
			logger.Init(Cfg.LogLevel)
			return nil
		},
	}
)

// -----------------------------------------------------------------------------
// Public entrypoint
// -----------------------------------------------------------------------------

// Execute is called by service/main.go.
func Execute() { _ = rootCmd.Execute() }

// -----------------------------------------------------------------------------
// CLI initialisation
// -----------------------------------------------------------------------------

func init() {
	// persistent flags available to *all* sub-commands
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"path to TOML/YAML/JSON config file")

	// register sub-commands (defined in their own files)
	rootCmd.AddCommand(
		serveCmd,
		migrateCmd,
		workerCmd,
		versionCmd,
	)
}
