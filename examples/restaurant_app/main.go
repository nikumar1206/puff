package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"reflect"

	"github.com/nikumar1206/puff"
	"github.com/nikumar1206/puff/middleware"
)

func main() {
	app := puff.DefaultApp("Restaurant Microservice")

	app.Use(middleware.Tracing())
	app.Use(middleware.CORS())
	app.Use(middleware.Panic())
	app.Use(middleware.Logging())

	g := puff.Field{
		PathParams: map[string]reflect.Kind{"name": reflect.String},
	}

	app.Get("/", g)

	f := puff.Field{
		PathParams: map[string]reflect.Kind{"name": reflect.String},
	}
	app.Get("/foos/{name}", f, func(c *puff.Context) {
		c.SendResponse(puff.GenericResponse{Content: "foo-bar"})
	})
	app.Get("/rawr", f, func(c *puff.Context) {
		c.SendResponse(puff.JSONResponse{Content: map[string]any{"foo": "bar"}, StatusCode: 200})
	})
	// app.IncludeRouter(routes.PizzaRouter())
	// dr := routes.DrinksRouter()
	// app.IncludeRouter(dr)

	app.SetDev()
	app.ListenAndServe()
}
