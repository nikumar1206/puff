package routes

import "github.com/nikumar1206/puff"

func WaterRouter() *puff.Router {
	r := puff.Router{Prefix: "/water"}
	r.Get("", "get water at no charge", func(puff.Request) interface{} {
		return puff.Response{
			Content: "dropping a bucket of water on you within 45 seconds",
		}
	})
	r.Post("", "add water to bucket", func(puff.Request) interface{} {
		return puff.Response{
			Content: "added water to bucket",
		}
	})
	return &r

}
