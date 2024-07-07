package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/nikumar1206/puff"
	"github.com/nikumar1206/puff/examples/restaurant_app/routes"
	"github.com/nikumar1206/puff/middleware"
)

func main() {
	app := puff.DefaultApp()

	fmt.Print(app.DocsURL)

	app.IncludeMiddlewares(
		middleware.CORS(), middleware.Panic(),
	)

	app.Get("/{$}", "Welcomes users to the application", func(ctx *puff.Context) {
		fmt.Print("did we get to here")
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
	app.IncludeRouter(routes.PizzaRouter())
	app.IncludeRouter(routes.DrinksRouter())
	app.ListenAndServe()
}
