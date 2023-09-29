package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/zer0go/ws-relay-service/internal/config"
	"os"
)

var rootCmd = &cobra.Command{
	Use:              "ws-relay",
	Short:            config.AppName,
	SilenceErrors:    true,
	PersistentPreRun: bootstrap,
}

func init() {
	rootCmd.
		PersistentFlags().
		CountP("verbosity", "v", "set logging verbosity")
}

func bootstrap(cmd *cobra.Command, _ []string) {
	verbosity, _ := cmd.Flags().GetCount("verbosity")
	config.ConfigureLogger(verbosity)

	err := config.Load()
	if err != nil {
		log.Warn().Err(err).Msg("config load failed")
	}
}

func Execute(version string) {
	rootCmd.Version = version
	rootCmd.Short += " " + version

	if err := rootCmd.Execute(); err != nil {
		log.Warn().Err(err).Msg("unexpected error")
		os.Exit(1)
	}
}
