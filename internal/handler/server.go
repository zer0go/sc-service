package handler

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/zer0go/ws-relay-service/internal/config"
	"net/http"
	"time"
)

type ServerHandler struct {
}

func NewServerHandler() *ServerHandler {
	return new(ServerHandler)
}

func (h *ServerHandler) Handle(cmd *cobra.Command, _ []string) error {
	addr, _ := cmd.Flags().GetString("addr")

	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Info().Msgf("%s listening on ws://%s", config.AppName, addr)

	return server.ListenAndServe()
}
