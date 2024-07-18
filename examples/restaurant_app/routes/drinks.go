package routes

import (
	"fmt"
	"time"

	"github.com/nikumar1206/puff"
)

func DrinksRouter() *puff.Router {
	r := puff.NewRouter(
		"Drinks",
		"/drinks",
	)

	r.Get("/stream-coca-cola", "", nil, func(c *puff.Context) {
		res := puff.StreamingResponse{
			StreamHandler: func(coca_cola *chan string) {
				for i := range 3 {
					*coca_cola <- fmt.Sprint(i)
					time.Sleep(time.Duration(2 * time.Second))
				}
			},
		}
		c.SendResponse(res)
	})
	r.Get("/{name}", "", nil, func(c *puff.Context) {
		res := puff.StreamingResponse{
			StreamHandler: func(coca_cola *chan string) {
				for i := range 3 {
					*coca_cola <- fmt.Sprint(i)
					time.Sleep(time.Duration(2 * time.Second))
				}
			},
		}
		c.SendResponse(res)
	})
	r.IncludeRouter(WaterRouter())
	r.IncludeRouter(SodaRouter())
	return r
}
