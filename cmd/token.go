package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zer0go/ws-relay-service/internal/handler"
)

var tokenCommand = &cobra.Command{
	Use:          "token",
	Short:        "Generate join token",
	SilenceUsage: true,
	RunE:         handler.NewTokenHandler().Handle,
}

func init() {
	rootCmd.AddCommand(tokenCommand)
}
