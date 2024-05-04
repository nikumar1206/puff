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
	return response.HTMLResponse{
		Content: "<h1>hello there from cheese</h1>",
	}
}

func main() {
	app := DefaultApp()

	app.RootRouter.GET(
		"/",

		"index route that says hello world",
		ex_rh,
	)

	v1Router := app.RootRouter.AddRouter(&router.Router{Prefix: "/v1"})

	v1Router.GET("/food", "hello there", ex2_rh)

	app.ListenAndServe()
}
