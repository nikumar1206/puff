package middleware

import (
	"fmt"
	"log/slog"

	"github.com/nikumar1206/puff"
	"github.com/nikumar1206/puff/utils"
)

func PanicMiddleware(next puff.HandlerFunc) puff.HandlerFunc {
	return func(c *puff.Context) {
		defer func() {
			a := recover()
			if a != nil {
				errorID := utils.RandomNanoID()
				slog.Error("Panic During Execution", slog.String("ERROR ID", errorID), slog.String("Error", a.(string)))
				errorMsg := fmt.Sprintf("There was a panic during the execution recovered by the panic handling middleware. Error ID: " + errorID)
				res := puff.GenericResponse{StatusCode: 500, Content: errorMsg}
				c.SendResponse(res)
			}
		}()
		next(c)
	}
}
