package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type loggerContextKey struct{}

// loggerFromContext falls back to slog.Default() so a handler or test that
// somehow runs outside the injectLogger middleware never panics on a nil
// logger — NewRouter always installs a non-nil one via injectLogger, this is
// pure defense in depth.
func loggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerContextKey{}).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

func injectLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), loggerContextKey{}, logger)))
		})
	}
}

// requestLoggerMiddleware logs one structured JSON line per request,
// correlated by request_id — replaces chi's middleware.Logger (colorized,
// unstructured stdout text meant for a human watching a dev terminal, not
// for a log aggregator) — architecture review finding #17. Must run after
// injectLogger (needs the logger from context) and before Recoverer's
// position in the old chain, i.e. outside it, so a panic's resulting 500
// still gets logged (matches chi's own Logger/Recoverer ordering contract).
func requestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		loggerFromContext(r.Context()).Info("http_request",
			"request_id", middleware.GetReqID(r.Context()),
			"method", r.Method,
			"path", r.URL.Path,
			"remote_ip", middleware.GetClientIP(r.Context()),
			"status", ww.Status(),
			"bytes", ww.BytesWritten(),
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}
