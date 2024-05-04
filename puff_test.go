package main

import (
	"testing"

	"github.com/nikumar1206/puff/request"
	"github.com/nikumar1206/puff/response"
	"github.com/nikumar1206/puff/router"
)

func example_route_handler(req request.Request) interface{} {
	return response.HTMLResponse{
		Content: "<h1>hello there</h1>",
	}
}

func TestApp(t *testing.T) {
	example_app := DefaultApp()

	example_router := router.Router{}

	example_router.GET(
		"",
		"index route that says hello world",
		example_route_handler,
	)
	example_app.RootRouter.IncludeRouter(&example_router)
	example_app.ListenAndServe()
}
