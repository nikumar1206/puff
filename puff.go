package puff

import (
	app "puff/app"
	"puff/logger"
)

func App(config app.Config) *app.App {
	if config.Port == 0 {
		config.Port = 8000
	}
	if config.Version == "" {
		config.Version = "1.0.0"
	}

	logger.DefaultLogger()
	return &app.App{Config: &config}
}

func DefaultApp() *app.App {
	logger.DefaultLogger()

	c := app.Config{
		Version:     "1.0.0",
		Name:        "Untitled",
		Network:     true,
		Port:        8000,
		OpenAPIDocs: true,
	}

	return App(c)
}
