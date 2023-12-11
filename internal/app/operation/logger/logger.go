package logger

import (
	"os"

	"github.com/kinneko-de/restaurant-file-store-svc/build"
	"github.com/rs/zerolog"
)

const LogLevel = zerolog.InfoLevel

var Logger zerolog.Logger = zerolog.New(os.Stdout).With().
	Timestamp().
	Caller().
	Str("version", build.Version).
	Logger()

func SetLogLevel(level zerolog.Level) {
	zerolog.SetGlobalLevel(level)
}
