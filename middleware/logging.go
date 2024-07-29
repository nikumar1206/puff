package middleware

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/nikumar1206/puff"
	color "github.com/nikumar1206/puff/color"
)

type LoggingConfig struct {
	LoggingFunction func(ctx puff.Context, startTime time.Time)
}

var DefaultLoggingConfig LoggingConfig = LoggingConfig{
	LoggingFunction: func(ctx puff.Context, startTime time.Time) {
		processingTime := time.Since(startTime).String()
		sc := ctx.GetStatusCode()
		var statusColor = fmt.Sprintf(" %d ", sc)
		switch {
		case sc >= 500:
			statusColor = color.ColorizeBold(string(sc), color.BgBrightRed, color.FgBlack)
		case sc >= 400:
			statusColor = color.ColorizeBold(string(sc), color.BgBrightYellow, color.FgBlack)
		case sc >= 300:
			statusColor = color.ColorizeBold(string(sc), color.BgBrightCyan, color.FgBlack)
		default:
			statusColor = color.ColorizeBold(string(sc), color.BgBrightGreen, color.FgBlack)
		}
		// TODO: make the below configurable
		// Request ID should only be present if present
		slog.Info(
			fmt.Sprintf("%s %s| %s | %s | %s ",
				statusColor,
				fmt.Sprintf("%s %s\t", ctx.Request.Method, ctx.Request.URL.String()),
				processingTime,
				fmt.Sprintf(ctx.GetRequestID()),
				fmt.Sprintf(ctx.ClientIP()),
			),
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
