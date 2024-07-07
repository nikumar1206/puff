package middleware

import (
	"github.com/google/uuid"
	"github.com/nikumar1206/puff"
)

func TracingMiddleware(next puff.HandlerFunc) puff.HandlerFunc {
	return func(c *puff.Context) {
		c.SetHeader("X-Request-ID", uuid.NewString())
		next(c)
	}
}
