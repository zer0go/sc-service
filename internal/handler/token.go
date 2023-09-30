package handler

import (
	"encoding/json"
	"github.com/mdp/qrterminal/v3"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/zer0go/ws-relay-service/internal/config"
	"os"
)

type TokenHandler struct {
	qrConfig qrterminal.Config
}

type qrCodeData struct {
	Url       string `json:"url"`
	JWTSecret string `json:"secret"`
}

func NewTokenHandler() *TokenHandler {
	return &TokenHandler{
		qrConfig: qrterminal.Config{
			Level:          qrterminal.L,
			Writer:         os.Stdout,
			HalfBlocks:     true,
			BlackChar:      qrterminal.BLACK_BLACK,
			WhiteBlackChar: qrterminal.WHITE_BLACK,
			WhiteChar:      qrterminal.WHITE_WHITE,
			BlackWhiteChar: qrterminal.BLACK_WHITE,
			QuietZone:      1,
		},
	}
}

func (h *TokenHandler) Handle(*cobra.Command, []string) error {
	data := qrCodeData{
		Url:       config.App.ServiceWebSocketUrl,
		JWTSecret: config.App.JWTSecret,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Err(err).Msg("encode to json failed")
		return nil
	}

	log.Debug().
		Interface("correctionLevel", h.qrConfig.Level).
		Int("quietZone", h.qrConfig.QuietZone).
		Bool("halfBlocks", h.qrConfig.HalfBlocks).
		Bytes("content", jsonData).
		Msg("generating qr code")

	qrterminal.GenerateWithConfig(string(jsonData), h.qrConfig)

	return nil
}
