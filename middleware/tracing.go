package middleware

import (
	"net/http"

	"github.com/google/uuid"
)

func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Request-ID", uuid.NewString())
		next.ServeHTTP(w, r)
	})
}
