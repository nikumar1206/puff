package main

import (
	"fmt"
	"github.com/nikumar1206/puff"
)

type HelloWorldInput struct {
	Name string `kind:"QueryParam" description:"Specify a name to say hello to."`
}

func main() {
	app := puff.DefaultApp()
	app.Get("/", "Hello, world!", func(c *puff.Context, input *struct {
		Name string `kind:"QueryParam" description:"Specify a name to say hello to."`
	}) {
		c.SendResponse(puff.GenericResponse{
			StatusCode: 200,
			Content:    fmt.Sprintf("Hello, %s!", input.Name),
		})
	})
}
