package routes

import (
	"github.com/nikumar1206/puff"
)

type OrderSodaInput struct {
	Name string `kind:"path" description:"soda to order"`
}

func SodaRouter() *puff.Router {
	r := puff.NewRouter("Soda", "/soda")
	r.Get("/", "", nil, func(c *puff.Context) {
		res := puff.GenericResponse{
			Content: "dropping a bucket of water on you within 45 seconds",
		}
		c.SendResponse(res)
	})
	r.Post("/order/{name}", "", nil, func(c *puff.Context) {
	})
	return r
}
