package main

import (
	"github.com/nikumar1206/puff/app"
	"github.com/nikumar1206/puff/router"

	"github.com/nikumar1206/puff/logger"
)

func NewApp(c *app.Config, r *router.Router) *app.App {
	if r == nil {
		r = &router.Router{
			Prefix: "",
		}
	}
	return &app.App{
		Config:     c,
		RootRouter: r,
	}
}

func DefaultApp() *app.App {
	logger.DefaultPuffLogger()

	c := app.Config{
		Network: true,
		Port:    8000,
	}

	r := router.Router{
		Prefix: "/api",
	}
	return NewApp(&c, &r)
}
