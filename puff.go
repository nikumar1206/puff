package main

import (
	"github.com/nikumar1206/puff/app"
	"github.com/nikumar1206/puff/logger"
	"github.com/nikumar1206/puff/router"
)

func App(c *app.Config) *app.App {
	r := &router.Router{
		Prefix: "",
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

	return App(&c)
}
