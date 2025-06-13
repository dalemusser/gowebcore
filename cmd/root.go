package cmd

import (
	"github.com/dalemusser/gowebcore/config"
	"github.com/dalemusser/gowebcore/logger"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "gowebsvc",
		Short: "gowebcore service runner",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := config.Load(&Cfg, config.WithConfigFile(cfgFile)); err != nil {
				return err
			}
			logger.Init(Cfg.LogLevel)
			return nil
		},
	}
	Cfg config.Base
)

func Execute() { _ = rootCmd.Execute() }

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.AddCommand(serveCmd, migrateCmd, workerCmd, versionCmd)
}
