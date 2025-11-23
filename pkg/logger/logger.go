package logger

import (
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func New(env string) *slog.Logger {

	var level slog.Level

	switch env {
	case envLocal, envDev:
		level = slog.LevelDebug
	case envProd:
		level = slog.LevelInfo
	default:
		level = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}))
}
