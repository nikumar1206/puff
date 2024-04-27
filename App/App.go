package app

import (
	"fmt"
	"net/http"
	handler "puff/Handler"
	route "puff/Route"
	router "puff/Router"
)

type Config struct {
	Network bool // host to the entire network?
	Reload  bool // live reload?
	Port    int  // port number to use
}

type App struct {
	*Config
	Routes  []route.Route
	Routers []*router.Router
	// Middlewares
}

func (a *App) IncludeRouter(r *router.Router) {
	a.Routers = append(a.Routers, r)
}

func (a *App) sendToHandler(w http.ResponseWriter, req *http.Request) {
	handler.Handler(w, req, a.Routers)
}

func (a *App) ListenAndServe() {
	http.HandleFunc("/", a.sendToHandler)
	addr := ""
	if a.Network {
		addr += "0.0.0.0"
	}
	addr += fmt.Sprintf(":%d", a.Port)

	fmt.Printf("Running app on port %d 🚀", a.Port)
	http.ListenAndServe(addr, nil)
}
