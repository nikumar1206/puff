package main

import app "puff/app"

func App(config app.Config) *app.App {
	return &app.App{Config: &config}
}

func DefaultApp() *app.App {
	// FIX_ME: reload bool should pick up from APP_ENV
	c := app.Config{
		Network: true,
		Reload:  true,
		Port:    8000,
	}

	return App(c)
}

func main() {
	app := DefaultApp()
	app.ListenAndServe()
}
