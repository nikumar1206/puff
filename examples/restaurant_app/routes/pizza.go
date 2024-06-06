package routes

import (
	"time"

	"github.com/nikumar1206/puff"
)

func getPizza(request puff.Request) interface{} {
	return puff.JSONResponse{
		Content: map[string]interface{}{
			"name": "Margherita Pizza",
			"ingredients": map[string]string{
				"Pizza dough":       "1 ball",
				"Tomato sauce":      "1/2 cup",
				"Mozzarella cheese": "1 cup, shredded",
				"Basil leaves":      "Handful",
				"Olive oil":         "1 tablespoon",
			},
			"instructions": []string{
				"Roll out dough.",
				"Spread tomato sauce.",
				"Add mozzarella cheese.",
				"Top with basil leaves.",
				"Drizzle olive oil.",
				"Bake at 475°F for 10-12 minutes.",
			},
		}}
}

func PizzaRouter() *puff.Router {
	r := &puff.Router{
		Name:   "Pizza related APIs for the restaurant",
		Prefix: "/pizza",
	}

	r.Get("", "Returns the greatest piza recipe you will ever find.", getPizza)

	r.Post("", "Places an order for a pizza.", func(req puff.Request) interface{} {
		timeOut := time.Duration(5)
		time.Sleep(timeOut)
		return puff.JSONResponse{
			StatusCode: 201,
			Content:    map[string]interface{}{"completed": true, "waitTime": timeOut},
		}
	})

	r.Patch("", "Unburns a burnt pizza.", func(puff.Request) interface{} {
		return puff.JSONResponse{
			StatusCode: 400,
			Content:    map[string]interface{}{"message": "Unburning a pizza is impossible."},
		}
	})
	ThumbnailFileResp := puff.FileResponse{FileName: "assets/chezpiza.jpg"}
	r.Get("/thumbnail", "returns thumbnail of piza", ThumbnailFileResp.Handler())

	return r
}
