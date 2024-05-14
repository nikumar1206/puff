package routes

import "github.com/nikumar1206/puff"

func DrinksRouter() *puff.Router {
	r := puff.Router{
		Name:   "All the drinks available at the store",
		Prefix: "/drinks",
	}
	r.IncludeRouter(WaterRouter())
	r.IncludeRouter(SodaRouter())
	return &r
}
