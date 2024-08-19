package main

import (
	"log/slog"
	"time"

	"github.com/nikumar1206/puff"
	"github.com/nikumar1206/puff/middleware"
)

type User struct {
	Foo string `json:"foo"`
}

func main() {
	app := puff.DefaultApp("Restaurant Microservice")
	app.DocsReload = true

	app.Use(middleware.Tracing())
	app.Use(middleware.CORS())
	app.Use(middleware.Logging())
	app.Use(middleware.CSRF())

	getRoute := app.Get("/", "", nil, func(c *puff.Context) {
		c.SendResponse(puff.FileResponse{
			FilePath: "examples/restaurant_app/assets/hello_world.html",
		})
	})
	getRoute.Description = "updating the description to foo"
	getRoute.Responses = map[int]puff.Response{
		200: puff.JSONResponse{StatusCode: 200, Content: &User{}},
	}

	// app.Get("/foos/{name}", "", nil, func(c *puff.Context) {
	// 	c.SendResponse(puff.GenericResponse{Content: "foo-bar"})
	// })
	// app.Get("/rawr", "", nil, func(c *puff.Context) {
	// 	c.SendResponse(puff.StreamingResponse{
	// 		StreamHandler: func(coca_cola *chan puff.ServerSideEvent) {
	// 			for i := range 3 {
	// 				*coca_cola <- puff.ServerSideEvent{Data: strconv.Itoa(i), Event: "foo", ID: puff.RandomNanoID(), Retry: 2}
	// 				time.Sleep(time.Duration(2 * time.Second))
	// 			}
	// 		}},
	// 	)
	// })

	// app.IncludeRouter(PastaRouter())
	// app.IncludeRouter(routes.PizzaRouter())
	// app.IncludeRouter(routes.DrinksRouter())
	// app.IncludeRouter(routes.SodaRouter())
	// app.IncludeRouter(routes.WaterRouter())

	app.Logger = puff.NewLogger(puff.LoggerConfig{
		Level:      slog.LevelDebug,
		Colorize:   true,
		TimeFormat: time.DateTime,
	})
	app.ListenAndServe(":8000")
}
