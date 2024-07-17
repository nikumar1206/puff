package middleware

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/fatih/color"
	"github.com/nikumar1206/puff"
)

type LoggingConfig struct {
	LoggingFunction func(ctx puff.Context, startTime time.Time)
}

var DefaultLoggingConfig LoggingConfig = LoggingConfig{
	LoggingFunction: func(ctx puff.Context, startTime time.Time) {
		processingTime := time.Since(startTime).String()
		sc := ctx.GetStatusCode()
		var statusColor *color.Color
		switch {
		case sc >= 500:
			statusColor = color.New(color.BgHiRed, color.FgBlack)
		case sc >= 400:
			statusColor = color.New(color.BgHiYellow, color.FgBlack)
		case sc >= 300:
			statusColor = color.New(color.BgHiCyan, color.FgBlack)
		default:
			statusColor = color.New(color.BgHiGreen, color.FgBlack)
		}
		// TODO: make the below configurable
		// Request ID should only be present if present
		slog.Info(
			fmt.Sprintf("|%s|%s|%s|%s|%s|",
				statusColor.Sprint(fmt.Sprintf(" %d ", sc)),
				fmt.Sprintf("\t%s %s\t", ctx.Request.Method, ctx.Request.URL.String()),
				processingTime,
				fmt.Sprintf(ctx.GetRequestID()),
				fmt.Sprintf(ctx.Request.RemoteAddr),
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
