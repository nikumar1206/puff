package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

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
	app.Get("/", g, func(ctx *puff.Context) {
		path, err := filepath.Abs("examples/restaurant_app/assets/hello_world.html")
		if err != nil {
			panic(err)
		}
		html_file, err := os.ReadFile(path)

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
	app.Post("/foos/{name}", f, func(c *puff.Context) {
		c.SendResponse(puff.GenericResponse{Content: "foo-bar"})
	})
	app.Get("/rawr", f, func(c *puff.Context) {
		c.SendResponse(
			puff.StreamingResponse{StreamHandler: func(coca_cola *chan puff.ServerSideEvent) {
				for i := range 3 {
					*coca_cola <- puff.ServerSideEvent{Data: strconv.Itoa(i), Event: "foo", ID: puff.RandomNanoID(), Retry: 2}
					time.Sleep(time.Duration(2 * time.Second))
				}
			}, StatusCode: 200},
		)
	})
	// app.IncludeRouter(routes.PizzaRouter())
	// dr := routes.DrinksRouter()
	// app.IncludeRouter(dr)

	app.Logger = puff.NewLogger(puff.LoggerConfig{
		UseJSON:   false,
		AddSource: false,
		Level:     slog.LevelDebug,
	})
	slog.Info("hello there")
	app.ListenAndServe(":8000")

}
