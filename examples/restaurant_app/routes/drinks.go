package routes

import (
	"github.com/nikumar1206/puff"
)

func DrinksRouter() *puff.Router {
	r := puff.NewRouter(
		"Drinks",
		"/drinks",
	)

	r.IncludeRouter(WaterRouter())
	r.IncludeRouter(SodaRouter())
	return r
}
