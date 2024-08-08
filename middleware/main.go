// Package middleware provides middlewares for handling common web application requirements.
package middleware

import "github.com/nikumar1206/puff"

// DefaultSkipper can be set on a middleware config to never skip the middleware
func DefaultSkipper(c *puff.Context) bool { return false }
