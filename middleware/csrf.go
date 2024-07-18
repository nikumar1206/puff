package middleware

import (
	"net/http"

	"github.com/nikumar1206/puff"
)

type CSRFMiddlewareConfig struct {
	CookieLength     int
	MaxAge           int
	ExpectedHeader   string   // ExpectedHeader declares what request header CSRF should be looking for when verifying.
	ProtectedMethods []string // ProtectedMethods declares what http methods CSRF should secure.
}

var DefaultCSRFMiddleware *CSRFMiddlewareConfig = &CSRFMiddlewareConfig{
	CookieLength:   32,
	MaxAge:         31449600,
	ExpectedHeader: "X-CSRFMiddlewareToken",
}

func createCSRFMiddleware(config *CSRFMiddlewareConfig) puff.Middleware {
	cookie_name := "CSRFMiddlewareToken"
	return func(next puff.HandlerFunc) puff.HandlerFunc {
		return func(c *puff.Context) {
			if c.Request.Method != "POST" {
				next(c)
				return
			}
			if c.GetCookie(cookie_name) != c.GetHeader(config.ExpectedHeader) {
				c.Forbidden("CSRFMiddlewareToken missing or incorrect.")
				return
			}
			c.SetCookie(&http.Cookie{
				Name:   cookie_name,
				Value:  puff.RandomToken(config.CookieLength),
				MaxAge: config.MaxAge, //expires after hour or session whichever comes first
			})
			next(c)
		}
	}
}

// CSRF middleware
func CSRF() puff.Middleware {
	return createCSRFMiddleware(DefaultCSRFMiddleware)
}

func CSRFWithConfig(config *CSRFMiddlewareConfig) puff.Middleware {
	return createCSRFMiddleware(config)
}
