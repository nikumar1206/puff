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

func PastaRouter() *puff.Router {
	pasta_router := puff.NewRouter("Pasta", "/pasta")

	pasta_home_input := new(PastaHomeInput)
	pasta_router.Get("/home", pasta_home_input, func(c *puff.Context) {
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

	return pasta_router
}
