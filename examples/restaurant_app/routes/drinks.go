package routes

import (
	"github.com/nikumar1206/puff"
)

type EchoInput struct {
	Body string
}

func DrinksRouter() *puff.Router {
	r := puff.NewRouter(
		"Drinks",
		"/drinks",
	)
	echoInput := new(EchoInput)

	r.Get("/echo", "echos the body", echoInput, func(c *puff.Context) {
		c.SendResponse(puff.GenericResponse{
			Content: echoInput.Body,
		})
	})

	return r
}
