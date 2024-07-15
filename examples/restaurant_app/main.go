package main

import (
	"log/slog"
	"os"
	"reflect"

	"github.com/nikumar1206/puff"
	"github.com/nikumar1206/puff/middleware"
)

func main() {
	app := puff.DefaultApp()

	app.IncludeMiddlewares(
		middleware.Tracing(), middleware.CORS(), middleware.Panic(), middleware.Logging(),
	)

	g := puff.Field{
		PathParams: map[string]reflect.Kind{"name": reflect.String},
	}
	app.Get("/", g, func(ctx *puff.Context) {
		html_file, err := os.ReadFile("assets/hello_world.html")

		switch {
		case err != nil:
			slog.Error(err.Error())
			res := puff.HTMLResponse{
				StatusCode: 500,
				Content:    "<h1>Sorry, an internal server error occured, and we couldn't read a file.</h1>",
			}
			ctx.SendResponse(res)

		case len(html_file) == 0:
			res := puff.HTMLResponse{
				StatusCode: 500,
				Content:    "<h1>Sorry, an internal server error occured, and reading a file gave us no bytes.</h1>",
			}
			ctx.SendResponse(res)
		default:
			res := puff.HTMLResponse{
				StatusCode: 200,
				Content:    string(html_file),
			}
			ctx.SendResponse(res)
		}

	})
	f := puff.Field{
		PathParams: map[string]reflect.Kind{"name": reflect.String},
	}
	app.Get("/foos/{name}", f, func(c *puff.Context) {
		c.SendResponse(puff.GenericResponse{Content: "foo-bar"})
	})
	// app.Get("/", f, func(c *puff.Context) {
	// 	c.SendResponse(puff.GenericResponse{Content: "foo-bar"})
	// })
	// app.IncludeRouter(routes.PizzaRouter())
	// dr := routes.DrinksRouter()
	// app.IncludeRouter(dr)

	app.ListenAndServe()
}
