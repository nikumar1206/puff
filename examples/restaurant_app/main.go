package main

import (
	"github.com/nikumar1206/puff"
	"github.com/nikumar1206/puff/examples/restaurant_app/routes"
	"github.com/nikumar1206/puff/middleware"
)

func main() {
	app := puff.DefaultApp("Restaurant Microservice")
	app.Config.DocsReload = true

	app.Use(middleware.Tracing())
	app.Use(middleware.CORS())
	app.Use(middleware.Panic())
	app.Use(middleware.Logging())
	app.Use(middleware.CSRF())
	app.Get("/", nil, func(c *puff.Context) {
		c.SendResponse(puff.FileResponse{
			FilePath: "examples/restaurant_app/assets/hello_world.html",
		})
	})

	app.Get("/foos/{name}", nil, func(c *puff.Context) {
		c.SendResponse(puff.GenericResponse{Content: "foo-bar"})
	})
	app.Get("/rawr", nil, func(c *puff.Context) {
		c.SendResponse(puff.JSONResponse{Content: map[string]any{"foo": "bar"}, StatusCode: 200})
	})

	app.IncludeRouter(PastaRouter())
	app.IncludeRouter(routes.PizzaRouter())
	app.IncludeRouter(routes.DrinksRouter())

	app.SetDev()
	app.ListenAndServe(":8000")
}
