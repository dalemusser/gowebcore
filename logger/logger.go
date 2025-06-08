package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

// lg is the global, structured logger returned by Instance().
var lg *slog.Logger

// Init must be called once (typically in main()).
// level accepts: debug | info | warn | error
func Init(level string) {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     parseLevel(level),
		AddSource: true, // adds file:line of the call site
	})

	// wrap the JSON handler so we can enrich every record with trace/span IDs.
	h := &otelHandler{next: jsonHandler}

	lg = slog.New(h)
}

// Instance returns the configured *slog.Logger.
func Instance() *slog.Logger { return lg }

/* -------------------------------------------------------------------------- */
/*                               helper types                                 */
/* -------------------------------------------------------------------------- */

// otelHandler adds {"trace_id": "...", "span_id": "..."} whenever a span
// is active in the record’s context, then delegates to the wrapped handler.
type otelHandler struct {
	next slog.Handler
}

func (h *otelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *otelHandler) Handle(ctx context.Context, r slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if span != nil && span.SpanContext().IsValid() {
		r.AddAttrs(
			slog.String("trace_id", span.SpanContext().TraceID().String()),
			slog.String("span_id", span.SpanContext().SpanID().String()),
		)
	}
	return h.next.Handle(ctx, r)
}

func (h *otelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &otelHandler{next: h.next.WithAttrs(attrs)}
}

func (h *otelHandler) WithGroup(name string) slog.Handler {
	return &otelHandler{next: h.next.WithGroup(name)}
}

/* -------------------------------------------------------------------------- */
/*                               level parser                                 */
/* -------------------------------------------------------------------------- */

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
