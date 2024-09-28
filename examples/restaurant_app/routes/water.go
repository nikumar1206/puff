package routes

import "github.com/nikumar1206/puff"

func WaterRouter() *puff.Router {
	r := puff.NewRouter("Water", "/water")

	// Retrieves water at no charge.
	r.Get("", nil, func(c *puff.Context) {
		res := puff.GenericResponse{
			Content: "dropping the bucket of water on you within 45 seconds",
		}
		c.SendResponse(res)
	})
	// Adds a water to the bucket.
	r.Post("", nil, func(c *puff.Context) {
		res := puff.GenericResponse{
			Content: "added water to bucket",
		}
		c.SendResponse(res)
	})
	return r
}
