package middleware

import (
	"log/slog"
	"time"

	"github.com/nikumar1206/puff"
)

type LoggingConfig struct {
	LoggingFunction func(ctx puff.Context, startTime time.Time)
}

var DefaultLoggingConfig LoggingConfig = LoggingConfig{
	LoggingFunction: func(ctx puff.Context, startTime time.Time) {
		processingTime := time.Since(startTime).String()
		slog.Info(
			"HTTP Request",
			slog.String("HTTP METHOD", ctx.Request.Method),
			slog.String("URL", ctx.Request.URL.String()),
			slog.String("Processing Time", processingTime),
		)
	},
}

func createLoggingMiddleware(lc LoggingConfig) puff.Middleware {
	return func(next puff.HandlerFunc) puff.HandlerFunc {
		return func(ctx *puff.Context) {
			startTime := time.Now()
			next(ctx)
			lc.LoggingFunction(*ctx, startTime)
		}
	}
}

func Logging() puff.Middleware {
	return createLoggingMiddleware(DefaultLoggingConfig)
}

func LoggingWithConfig(tc LoggingConfig) puff.Middleware {
	return createLoggingMiddleware(tc)
}
