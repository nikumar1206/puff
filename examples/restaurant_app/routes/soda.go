package routes

import (
	"github.com/nikumar1206/puff"
)

func SodaRouter() *puff.Router {
	r := puff.Router{Prefix: "/soda"}
	r.Get("", "get water for a dollar", func(c *puff.Context) {
		res := puff.GenericResponse{
			Content: "dropping a bucket of water on you within 45 seconds",
		}
		c.SendResponse(res)
	})
	r.Get("/fanta", "request fanta", func(c *puff.Context) {
		panic("WOAH! WE DONT SELL FANTA!!")
	})
	return &r
}
