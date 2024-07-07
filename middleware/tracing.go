package middleware

import (
	"github.com/google/uuid"
	"github.com/nikumar1206/puff"
)

type TracingConfig struct {
	TracerName  string
	IDGenerator func() string
}

var DefaultTracingConfig TracingConfig = TracingConfig{
	TracerName:  "X-Request-ID",
	IDGenerator: uuid.NewString,
}

func createTracingMiddleware(tc TracingConfig) puff.Middleware {
	return func(next puff.HandlerFunc) puff.HandlerFunc {
		return func(c *puff.Context) {
			c.SetHeader(tc.TracerName, tc.IDGenerator())
			next(c)
		}
	}
}

func Tracing() puff.Middleware {
	return createTracingMiddleware(DefaultTracingConfig)
}

func TracingWithConfig(tc TracingConfig) puff.Middleware {
	return createTracingMiddleware(tc)
}
