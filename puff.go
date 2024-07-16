package puff

import (
	"github.com/nikumar1206/puff/logger"
)

type HandlerFunc func(c *Context)
type Middleware func(next HandlerFunc) HandlerFunc

func App(c *Config) *PuffApp {
	r := &Router{Name: "Puff Default", Tag: "Default", Description: "Puff Default Router"}
	if c.Version == "" {
		c.Version = "0.0.0"
	}

	return &PuffApp{
		Config:     c,
		RootRouter: r,
	}
}

func DefaultApp(name string) *PuffApp {
	c := Config{
		Version:    "1.0.0",
		Name:       name,
		ListenAddr: ":8000",
		DocsURL:    "/docs",
	}
	a := App(&c)
	a.Logger = logger.DefaultPuffLogger()
	return a
}
