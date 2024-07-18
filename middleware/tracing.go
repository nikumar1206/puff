package middleware

import (
	"github.com/google/uuid"
	"github.com/nikumar1206/puff"
)

type TracingConfig struct {
	TracerName  string        //TracerName is the name of the response header where the request id will be.
	IDGenerator func() string // IDGenerator is a function that returns a string that will be used as the request id.
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

// Tracing middleware provides the ability to automatically trace every route with a request id.
func Tracing() puff.Middleware {
	return createTracingMiddleware(DefaultTracingConfig)
}

func TracingWithConfig(tc TracingConfig) puff.Middleware {
	return createTracingMiddleware(tc)
}
