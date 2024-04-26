package puff

import app "puff/App"

type AppConfig struct {
	Network bool
	Reload  bool
	Port    int
}

func App(config AppConfig) app.AppI {
	return app.AppI{
		Network: config.Network,
		Reload:  config.Reload,
		Port:    config.Port,
	}
}
