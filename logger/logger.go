package logger

import (
	"log/slog"
	"os"
	"strings"
)

var lg *slog.Logger

// Init sets the global structured logger.
func Init(level string) {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLevel(level),
		AddSource: true,
	})
	lg = slog.New(h)
}

func Instance() *slog.Logger { return lg }

func parseLevel(l string) slog.Level {
	switch strings.ToLower(l) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
