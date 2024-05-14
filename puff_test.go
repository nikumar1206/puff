package puff

import (
	"testing"
)

func example_route_handler(req Request) interface{} {
	return HTMLResponse{
		Content: "<h1>hello there</h1>",
	}
}

func TestApp(t *testing.T) {
	example_app := DefaultApp()

	example_router := Router{}

	example_router.Get(
		"",
		"index route that says hello world",
		example_route_handler,
	)
	example_app.RootRouter.IncludeRouter(&example_router)
	example_app.ListenAndServe()
}
