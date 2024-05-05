package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		requestID := uuid.NewString()
		next.ServeHTTP(w, r)
		processingTime := time.Since(startTime).String()
		w.Header().Add("X-Processing-Time", processingTime)
		w.Header().Add("X-Request-ID", requestID)
	})
}
