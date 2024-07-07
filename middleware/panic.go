package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/nikumar1206/puff"
	"github.com/nikumar1206/puff/utils"
)

type PanicConfig struct {
	FormatErrorResponse func(any) puff.Response
}

var DefaultPanicConfig PanicConfig = PanicConfig{
	FormatErrorResponse: func(a any) puff.Response {
		errorID := utils.RandomNanoID()
		slog.Error("Panic During Execution", slog.String("ERROR ID", errorID), slog.String("Error", a.(string)))
		errorMsg := fmt.Sprintf("There was a panic during the execution recovered by the panic handling middleware. Error ID: " + errorID)
		return puff.GenericResponse{StatusCode: http.StatusInternalServerError, Content: errorMsg}
	},
}

func createPanicMiddleware(pc PanicConfig) puff.Middleware {
	return func(next puff.HandlerFunc) puff.HandlerFunc {
		return func(c *puff.Context) {
			defer func() {
				a := recover()
				if a != nil {
					res := pc.FormatErrorResponse(a)
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
