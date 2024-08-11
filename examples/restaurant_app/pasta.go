package main

import (
	"fmt"
	"net/http"

	"github.com/nikumar1206/puff"
)

type PastaHomeInput struct {
	ID              int    `description:"proposed id of dish" kind:"query"`
	Name            string `description:"name of pasta dish" kind:"query"`
	LastDishOrdered string `description:"last dish ordered" kind:"cookie"`
}

type Pizza struct {
	Name string `json:"name"`
}
type PastaNewCheeseInput struct {
	HelloWorld Pizza `kind:"header"`
}

var cheeses = map[int]string{
	0: "Mozzerella",
	1: "Swiss",
}

type FooResponse struct {
	error   string
	message string
}

type PastaCheeseInput struct {
	Id int `kind:"path" description:"id of cheese"`
}

func PastaRouter() *puff.Router {
	pastaRouter := puff.NewRouter("Pasta", "/pasta")

	pasta_home_input := new(PastaHomeInput)
	pastaRouter.Get("/home", "", pasta_home_input, func(c *puff.Context) {
		if pasta_home_input.Name == pasta_home_input.LastDishOrdered {
			c.SendResponse(puff.GenericResponse{
				Content: fmt.Sprintf(
					"You've already ordered %s (%d). Don't worry it'll be ready soon!",
					pasta_home_input.LastDishOrdered,
					pasta_home_input.ID,
				),
			})
			return
		}
		c.SetCookie(&http.Cookie{
			Name:  "LastDishOrdered",
			Value: pasta_home_input.Name,
		})
		c.SendResponse(puff.GenericResponse{
			Content: fmt.Sprintf(
				"Making pasta dish '%s' with id %d. Check back in thirty minutes, it should be ready!",
				pasta_home_input.Name,
				pasta_home_input.ID,
			),
		})
	})

	pasta_cheese_input := new(PastaCheeseInput)
	pastaRouter.Get("/cheese/{Id}", "", pasta_cheese_input, func(c *puff.Context) {
		cheese, ok := cheeses[pasta_cheese_input.Id]
		if !ok {
			c.NotFound("Cheese with id %d not found.", pasta_cheese_input.Id)
			return
		}
		c.SendResponse(puff.GenericResponse{
			Content: cheese,
		})
	})

	pasta_newcheese_input := new(PastaNewCheeseInput)
	pastaRouter.Post("/cheese/{id}", "", pasta_newcheese_input, func(c *puff.Context) {
	})

	return pastaRouter
}
