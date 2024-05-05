package main

import (
	"github.com/nikumar1206/puff/request"
	"github.com/nikumar1206/puff/response"
	"github.com/nikumar1206/puff/router"
)

func ex_rh(req request.Request) interface{} {
	return response.HTMLResponse{
		Content: "<h1>hello there</h1>",
	}
}

func ex2_rh(req request.Request) interface{} {
	return response.JSONResponse{
		Content: map[string]interface{}{"hello there": "cheese;", "bloop": "scoop"},
	}
}

func main() {
	app := DefaultApp()

	//Food
	example_food := router.Router{
		Name:   "food",
		Prefix: "/food",
	}

	example_food.GET("/pizza", "Returns the greatest piza reciepe you will ever find.", func(req request.Request) interface{} {
		return response.HTMLResponse{
			Content: "<h1>pizza</h1>",
		}
	})
	example_food.GET("/pasta", "Returns the greatest pasta reciepe you will ever find.", func(req request.Request) interface{} {
		return response.HTMLResponse{
			StatusCode: 418,
			Content:    "<h1>no pasta for you</h1>",
		}
	})
	example_food.POST("/pizza", "Makes a pizza.", func(req request.Request) interface{} {
		return response.JSONResponse{
			StatusCode: 201,
			Content:    map[string]interface{}{"completed": true, "waitTime": 214},
		}
	})
	example_food.PATCH("/pizza", "Unburns a burnt pizza.", func(request.Request) interface{} {
		return response.Response{
			StatusCode: 400,
			Content:    "Unburning a pizza is impossible.",
		}
	})

	//Cheese
	example_cheese := router.Router{
		Name:   "cheese",
		Prefix: "/cheese",
	}

	example_cheese.POST("/gouda", "Request a wheel of gouda.", func(request.Request) interface{} {
		return response.JSONResponse{
			Content: map[string]interface{}{"completed": "yes", "charged": 100.36, "should_have_used_american": false},
		}
	})
	example_cheese.PUT("/swiss", "Puts a slice of swiss cheese on your dish.", func(request.Request) interface{} {
		return response.Response{}
	})
	//Drinks
	example_drinks := router.Router{
		Name:   "drinks",
		Prefix: "/drinks",
	}

	example_drinks.GET("/water", "get water at no charge", func(request.Request) interface{} {
		return response.Response{
			Content: "dropping a bucket of water on you within 45 seconds",
		}
	})
	example_drinks.POST("/water", "add water to bucket", func(request.Request) interface{} {
		return response.Response{
			Content: "added water to bucket",
		}
	})
	example_food.IncludeRouter(&example_cheese)
	app.IncludeRouter(&example_food)
	app.IncludeRouter(&example_drinks)
	app.ListenAndServe()
}
