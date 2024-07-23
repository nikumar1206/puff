package routes

import (
	"strings"
	"time"

	"github.com/nikumar1206/puff"
)

func getPizza(c *puff.Context) {
	res := puff.JSONResponse{
		Content: map[string]interface{}{
			"name": "Margherita Pizza",
			"ingredients": map[string]string{
				"Pizza dough":       "1 ball",
				"Tomato sauce":      "1/2 cup",
				"Mozzarella cheese": "1 cup, shredded",
				"Basil leaves":      "Handful",
				"Olive oil":         "1 tablespoon",
			},
			"instructions": []string{
				"Roll out dough.",
				"Spread tomato sauce.",
				"Add mozzarella cheese.",
				"Top with basil leaves.",
				"Drizzle olive oil.",
				"Bake at 475°F for 10-12 minutes.",
			},
		}}
	c.SendResponse(res)
}

type Pizza struct {
	Name        string
	Ingredients []string
}
type NewPizzaInput struct {
	// Pizza []string `kind:"body"`
	Pizza Pizza `kind:"body"`
}

func PizzaRouter() *puff.Router {
	r := &puff.Router{
		Name:   "Pizza",
		Prefix: "/pizza",
	}

	r.Get("/", "", nil, getPizza)

	r.Post("/", "", nil, func(c *puff.Context) {
		timeOut := 5 * time.Second
		time.Sleep(timeOut)
		res := puff.JSONResponse{
			StatusCode: 201,
			Content:    map[string]interface{}{"completed": true, "waitTime": timeOut},
		}
		c.SendResponse(res)
	})

	r.Patch("/", "", nil, func(c *puff.Context) {
		res := puff.JSONResponse{
			StatusCode: 400,
			Content:    map[string]interface{}{"message": "Unburning a pizza is impossible."},
		}
		c.SendResponse(res)
	})

	r.Get("/thumbnail", "", nil, func(c *puff.Context) {
		c.SendResponse(puff.FileResponse{
			FilePath: "examples/restaurant_app/assets/chezpiza.jpg",
		})
	})

	newPizzaInput := new(NewPizzaInput)
	r.Post("/new", "", newPizzaInput, func(c *puff.Context) {
		c.SendResponse(puff.GenericResponse{
			Content: "creating " + newPizzaInput.Pizza.Name + " with ingredients " + strings.Join(newPizzaInput.Pizza.Ingredients, ","),
		})
	})
	return r
}
