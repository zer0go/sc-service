package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zer0go/ws-relay-service/internal/handler"
)

var serverCommand = &cobra.Command{
	Use:          "server",
	Short:        "Start server",
	SilenceUsage: true,
	RunE:         handler.NewServerHandler().Handle,
}

func init() {
	serverCommand.Flags().StringP("addr", "a", "0.0.0.0:8080", "listening address")

	rootCmd.AddCommand(serverCommand)
}
