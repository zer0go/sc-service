package config

import (
	"github.com/caarlos0/env/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	AppName = "WebSocket Relay Server"
)

type config struct {
	JWTSecret           string `env:"JWT_SECRET,required"`
	ServiceWebSocketUrl string `env:"SERVICE_WS_URL,required"`
}

var App config

func Load() error {
	App = config{}

	return env.Parse(&App)
}

func ConfigureLogger(verbosity int) {
	switch verbosity {
	case 1:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case 0:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		dir := filepath.Dir(file)
		parent := filepath.Base(dir)
		return parent + "/" + filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	log.Logger = log.
		With().
		Stack().
		Caller().
		Logger()

	if os.Getenv("LOG_FORMAT") != "json" {
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: zerolog.TimeFieldFormat,
		})
	}
}
