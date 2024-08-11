package middleware

import (
	"github.com/google/uuid"
	"github.com/nikumar1206/puff"
)

// TracingConfig is a struct to configure the tracing middleware.
type TracingConfig struct {
	// Skip allows skipping the middleware for specific requests.
	// The function receives the request context and should return true if the middleware should be skipped.
	Skip func(*puff.Context) bool
	//TracerName is the name of the response header in which Request ID will be present.
	TracerName string
	// IDGenerator is a function that must return a string to generate the Request ID.
	IDGenerator func() string
}

// DefaultTracingConfig is a TracingConfig with specified default values.
var DefaultTracingConfig TracingConfig = TracingConfig{
	TracerName:  "X-Request-ID",
	IDGenerator: uuid.NewString,
	Skip:        DefaultSkipper,
}

// createCSRFMiddleware is used to create a CSRF middleware with a config.
func createTracingMiddleware(tc TracingConfig) puff.Middleware {
	return func(next puff.HandlerFunc) puff.HandlerFunc {
		return func(c *puff.Context) {
			if tc.Skip != nil && tc.Skip(c) {
				next(c)
				return
			}
			c.SetHeader(tc.TracerName, tc.IDGenerator())
			next(c)
		}
	}
}

// Tracing middleware provides the ability to automatically trace every route with a request id.
// The function returns a middleware with the default tracing config.
func Tracing() puff.Middleware {
	return createTracingMiddleware(DefaultTracingConfig)
}

// TracingWithConfig returns a tracing middleware with the config given.
func TracingWithConfig(tc TracingConfig) puff.Middleware {
	return createTracingMiddleware(tc)
}
