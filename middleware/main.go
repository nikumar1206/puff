// Package middleware provides middlewares for handling common web application requirements.
package middleware

import "github.com/ThePuffProject/puff"

// DefaultSkipper can be set on a middleware config to never skip the middleware
func DefaultSkipper(c *puff.Context) bool { return false }
