package routes

import "github.com/nikumar1206/puff"

func WaterRouter() *puff.Router {
	r := puff.NewRouter("Water", "/water")
	r.Get("", "get water at no charge", nil, func(c *puff.Context) {
		res := puff.GenericResponse{
			Content: "dropping a bucket of water on you within 45 seconds",
		}
		c.SendResponse(res)
	})
	r.Post("", "add water to bucket", nil, func(c *puff.Context) {
		res := puff.GenericResponse{
			Content: "added water to bucket",
		}
		c.SendResponse(res)
	})
	return r
}
