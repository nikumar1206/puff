package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/nikumar1206/puff"
)

type PanicConfig struct {
	FormatErrorResponse func(c puff.Context, err any) puff.Response
}

var DefaultPanicConfig PanicConfig = PanicConfig{
	FormatErrorResponse: func(c puff.Context, err any) puff.Response {
		errorID := puff.RandomNanoID()
		slog.Error("Panic During Execution", slog.String("ERROR ID", errorID), slog.Any("Error", err))
		errorMsg := fmt.Sprintf("There was a panic during the execution recovered by the panic handling middleware. Error ID: " + errorID)
		return puff.JSONResponse{StatusCode: http.StatusInternalServerError, Content: map[string]any{"error": errorMsg, "Request-ID": c.GetRequestID()}}
	},
}

func createPanicMiddleware(pc PanicConfig) puff.Middleware {
	return func(next puff.HandlerFunc) puff.HandlerFunc {
		return func(c *puff.Context) {
			defer func() {
				a := recover()
				if a != nil {
					res := pc.FormatErrorResponse(*c, a)
					c.SendResponse(res)
				}
			}()
			next(c)
		}
	}
}

func Panic() puff.Middleware {
	return createPanicMiddleware(DefaultPanicConfig)
}

func PanicWithConfig(pc PanicConfig) puff.Middleware {
	return createPanicMiddleware(pc)
}
