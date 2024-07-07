package middleware

import (
	"github.com/nikumar1206/puff"
)

type CORSMiddlewareConfig struct {
	AllowedOrigin  string
	AllowedMethods []string
	AllowedHeaders []string
}

func CORSMiddleware(next puff.HandlerFunc) puff.HandlerFunc {
	return func(ctx *puff.Context) {
		ctx.SetHeader("Access-Control-Allow-Origin", "*")
		ctx.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")
		next(ctx)
	}
}
