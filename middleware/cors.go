package middleware

import (
	"net/http"
	"strings"

	"github.com/nikumar1206/puff"
)

type CORSConfig struct {
	AllowedOrigin  string
	AllowedMethods []string
	AllowedHeaders []string
}

var DefaultCORSConfig CORSConfig = CORSConfig{
	AllowedOrigin: "*",
	AllowedMethods: []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
		http.MethodPut,
		http.MethodPatch,
		http.MethodPost,
		http.MethodDelete,
	},
	AllowedHeaders: []string{"Content-Type", "Authorization", "Origin"},
}

func createCORSMiddleware(c CORSConfig) puff.Middleware {
	allowedMethods := strings.Join(c.AllowedMethods, ",")
	allowedHeaders := strings.Join(c.AllowedHeaders, ",")

	return func(next puff.HandlerFunc) puff.HandlerFunc {
		return func(ctx *puff.Context) {
			ctx.SetHeader("Access-Control-Allow-Origin", c.AllowedOrigin)
			ctx.SetHeader("Access-Control-Allow-Methods", allowedMethods)
			ctx.SetHeader("Access-Control-Allow-Headers", allowedHeaders)
			next(ctx)
		}
	}
}

func CORS() puff.Middleware {
	return createCORSMiddleware(DefaultCORSConfig)
}

func CORSWithConfig(c CORSConfig) puff.Middleware {
	return createCORSMiddleware(c)
}
