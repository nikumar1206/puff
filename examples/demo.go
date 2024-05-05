package main

import (
	"github.com/nikumar1206/puff"
	"github.com/nikumar1206/puff/middleware"
)

func ex_rh(req puff.Request) interface{} {
	return puff.HTMLResponse{
		Content: "<h1>hello there</h1>",
	}
}

func ex2_rh(req puff.Request) interface{} {
	return puff.JSONResponse{
		Content: map[string]interface{}{"hello there": "cheese;", "bloop": "scoop"},
	}
}

func main() {
	app := puff.DefaultApp()

	app.Middlewares = []middleware.Middleware{
		middleware.TracingMiddleware,
		middleware.LoggingMiddleware,
		middleware.CORSMiddleware,
	}

	//Food
	example_food := puff.Router{
		Name:   "food",
		Prefix: "/food",
	}

	example_food.GET("/pizza", "Returns the greatest piza reciepe you will ever find.", func(req puff.Request) interface{} {
		return puff.HTMLResponse{
			Content: "<h1>pizza</h1>",
		}
	})
	example_food.GET("/pasta", "Returns the greatest pasta reciepe you will ever find.", func(req puff.Request) interface{} {
		return puff.HTMLResponse{
			StatusCode: 418,
			Content:    "<h1>no pasta for you</h1>",
		}
	})
	example_food.POST("/pizza", "Makes a pizza.", func(req puff.Request) interface{} {
		return puff.JSONResponse{
			StatusCode: 201,
			Content:    map[string]interface{}{"completed": true, "waitTime": 214},
		}
	})
	example_food.PATCH("/pizza", "Unburns a burnt pizza.", func(puff.Request) interface{} {
		return puff.Response{
			StatusCode: 400,
			Content:    "Unburning a pizza is impossible.",
		}
	})

	//Cheese
	example_cheese := puff.Router{
		Name:   "cheese",
		Prefix: "/cheese",
	}

	example_cheese.POST("/gouda", "puff.Request a wheel of gouda.", func(puff.Request) interface{} {
		return puff.JSONResponse{
			Content: map[string]interface{}{"completed": "yes", "charged": 100.36, "should_have_used_american": false},
		}
	})
	example_cheese.PUT("/swiss", "Puts a slice of swiss cheese on your dish.", func(puff.Request) interface{} {
		return puff.Response{}
	})
	//Drinks
	example_drinks := puff.Router{
		Name:   "drinks",
		Prefix: "/drinks",
	}

	example_drinks.GET("/water", "get water at no charge", func(puff.Request) interface{} {
		return puff.Response{
			Content: "dropping a bucket of water on you within 45 seconds",
		}
	})
	example_drinks.POST("/water", "add water to bucket", func(puff.Request) interface{} {
		return puff.Response{
			Content: "added water to bucket",
		}
	})
	example_food.IncludeRouter(&example_cheese)
	app.IncludeRouter(&example_food)
	app.IncludeRouter(&example_drinks)
	app.ListenAndServe()
}
