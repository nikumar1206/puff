package middleware

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
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

func RandomLogID() string {
	id := ""
	for i := range 8 {
		if i == 4 {
			id += "-"
		}
		if i < 4 {
			r := rand.IntN(25) + 1
			id += fmt.Sprintf("%c", ('A' - 1 + r))
		} else {
			r := rand.IntN(9)
			id += fmt.Sprint(r)
		}
	}
	return id
}
