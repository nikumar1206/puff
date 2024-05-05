package puff

import (
	"github.com/nikumar1206/puff/logger"
)

func App(c *Config) *PuffApp {
	r := &Router{
		Prefix: "",
	}
	if c.Version == "" {
		c.Version = "1.0.0"
	}

	return &PuffApp{
		Config:     c,
		RootRouter: r,
	}
}

func DefaultApp() *PuffApp {
	logger.DefaultPuffLogger()

	c := Config{
		Version: "1.0.0",
		Name:    "Untitled",
		Network: true,
		Port:    8000,
		DocsURL: "/docs",
	}

	return App(&c)
}
