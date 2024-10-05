package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/nikumar1206/puff"
	"github.com/nikumar1206/puff/middleware"
)

type User struct {
	Foo       string  `json:"foo"`
	Coolbeans int     `json:"boo"`
	Foobar    float32 `json:"troo"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func main() {
	app := puff.DefaultApp("Restaurant Microservice")
	app.DocsReload = true

	app.Use(middleware.Tracing())
	app.Use(middleware.CORS())
	app.Use(middleware.Logging())
	app.Use(middleware.CSRF())

	app.Get("/", "", nil, func(c *puff.Context) {
		c.SendResponse(puff.FileResponse{
			FilePath: "examples/restaurant_app/assets/hello_world.html",
		})
	}).WithResponse(http.StatusOK, puff.ResponseT[User])

	app.Get("/foo", "", nil, func(c *puff.Context) {
		c.SendResponse(puff.FileResponse{
			FilePath: "examples/restaurant_app/assets/hello_world.html",
		})
	}).WithResponses(
		puff.DefineResponse(http.StatusOK, puff.ResponseT[User]),
		puff.DefineResponse(http.StatusBadRequest, puff.ResponseT[ErrorResponse]),
		puff.DefineResponse(http.StatusInternalServerError, puff.ResponseT[ErrorResponse]),
	)

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
