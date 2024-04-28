package main

import (
	app "puff/app"
	"puff/logger"
)

func App(config app.Config) *app.App {
	if config.Port == 0 {
		config.Port = 8000
	}

	logger.DefaultLogger()
	return &app.App{Config: &config}
}

func DefaultApp() *app.App {
	logger.DefaultLogger()

	c := app.Config{
		Network: true,
		Port:    8000,
	}

	return App(c)
}
