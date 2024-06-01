package main

import (
	"log/slog"
	"os"

	"github.com/nikumar1206/puff"
	"github.com/nikumar1206/puff/examples/restaurant_app/routes"
	"github.com/nikumar1206/puff/middleware"
)

func main() {
	app := puff.DefaultApp()

	app.IncludeMiddlewares(
		middleware.CORSMiddleware,
	)

	app.Get("/{$}", "Welcomes users to the application", func(req puff.Request) puff.Response {
		html_file, err := os.ReadFile("assets/hello_world.html")

		switch {
		case err != nil:
			slog.Error(err.Error())
			return puff.HTMLResponse{
				StatusCode: 500,
				Content:    "<h1>Sorry, an internal server error occured, and we couldn't read a file.</h1>",
			}

		case len(html_file) == 0:
			return puff.HTMLResponse{
				StatusCode: 500,
				Content:    "<h1>Sorry, an internal server error occured, and reading a file gave us no bytes.</h1>",
			}
		default:
			return puff.HTMLResponse{
				StatusCode: 200,
				Content:    string(html_file),
			}
		}
	})
	app.IncludeRouter(routes.PizzaRouter())
	app.IncludeRouter(routes.DrinksRouter())
	app.ListenAndServe()
}
