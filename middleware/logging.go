package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		processingTime := time.Since(startTime).String()
		slog.Info(
			"HTTP Request",
			slog.String("HTTP METHOD", r.Method),
			slog.String("URL", r.URL.String()),
			slog.String("Processing Time", processingTime),
		)
	})
}
