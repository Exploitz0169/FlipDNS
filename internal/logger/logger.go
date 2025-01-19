package logger

import (
	"log/slog"
	"os"
)

func InitLogger() {

	level := slog.LevelDebug
	if os.Getenv("ENVIRONMENT") == "production" {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	slog.SetDefault(logger)

}
