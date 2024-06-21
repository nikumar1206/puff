package middleware

import (
	"log/slog"
	"time"

	"github.com/nikumar1206/puff"
)

func LoggingMiddleware(next puff.HandlerFunc) puff.HandlerFunc {
	return func(ctx *puff.Context) {
		startTime := time.Now()
		next(ctx)
		processingTime := time.Since(startTime).String()
		slog.Info(
			"HTTP Request",
			slog.String("HTTP METHOD", ctx.Request.Method),
			slog.String("URL", ctx.Request.URL.String()),
			slog.String("Processing Time", processingTime),
		)
	}
}
