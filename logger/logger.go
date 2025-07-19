// gowebcore/logger/logger.go
package logger

import (
	"bufio"
	"io"
	"log/slog"
	"os"
	"strings"
)

// lg is the global slog.Logger configured by Init().
var lg *slog.Logger

// Init must be called once (e.g. from main.go).
// level = "debug" | "info" | "warn" | "error"
func Init(level string) {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	handlerOpts := &slog.HandlerOptions{
		Level:     lvl,
		AddSource: true, // file:line in each log entry
	}

	handler := slog.NewJSONHandler(os.Stderr, handlerOpts)
	lg = slog.New(handler)

	// Make slog.* helpers use the same logger.
	slog.SetDefault(lg)
}

/*────────────────── convenience helpers ──────────────────*/

// L returns the global logger (advanced use).
func L() *slog.Logger { return lg }

func Debug(msg string, attrs ...any) { lg.Debug(msg, attrs...) }
func Info(msg string, attrs ...any)  { lg.Info(msg, attrs...) }
func Warn(msg string, attrs ...any)  { lg.Warn(msg, attrs...) }
func Error(msg string, attrs ...any) { lg.Error(msg, attrs...) }

/*────────────────── io.Writer shim (optional) ─────────────*/

// writer logs each Write() call at Info level.
// Useful when a library insists on an io.Writer for logs.
type writer struct {
	l *slog.Logger
}

func (w writer) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(p)))
	for scanner.Scan() {
		w.l.Info(scanner.Text())
	}
	return len(p), scanner.Err()
}

// Writer returns an io.Writer that logs at Info level.
func Writer() io.Writer { return writer{l: lg} }
