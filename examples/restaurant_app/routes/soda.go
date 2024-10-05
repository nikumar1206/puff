package routes

import (
	"github.com/nikumar1206/puff"
	"github.com/nikumar1206/puff/middleware"
)

type OrderSodaInput struct {
	Name string `kind:"path" description:"soda to order"`
}

func SodaRouter() *puff.Router {
	r := puff.NewRouter("Soda", "/soda")
	r.Use(middleware.Panic())
	r.Get("/", nil, func(c *puff.Context) {
		res := puff.GenericResponse{
			Content: "dropping a bucket of water on you within 45 seconds",
		}
		c.SendResponse(res)
	})
	// Retrieves fanta.
	r.Get("/fanta", nil, func(c *puff.Context) {
		panic("we no serve fanta.")
	})
	return r
}
