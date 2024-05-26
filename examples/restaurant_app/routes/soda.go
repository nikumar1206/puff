package routes

import (
	"github.com/nikumar1206/puff"
)

func SodaRouter() *puff.Router {
	r := puff.Router{Prefix: "/soda"}
	r.Get("", "get water for a dollar", func(puff.Request) interface{} {
		return puff.Response{
			Content: "dropping a bucket of water on you within 45 seconds",
		}
	})
	r.Get("/fanta", "request fanta", func(r puff.Request) interface{} {
		panic("WOAH! WE DONT SELL FANTA!!")
	})
	return &r
}
