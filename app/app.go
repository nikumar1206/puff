package app

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/nikumar1206/puff/handler"
	"github.com/nikumar1206/puff/middleware"
	"github.com/nikumar1206/puff/openapi"
	"github.com/nikumar1206/puff/route"
	"github.com/nikumar1206/puff/router"
)

type Config struct {
	Network bool // host to the entire network?
	Port    int  // port number to use
	OpenAPI *openapi.OpenAPI
}

type App struct {
	*Config
	RootRouter *router.Router // This is the root router. All other routers will work underneath this.
	// add middlewares
}

// gets all routes for a router
func (a *App) GetRoutes(r *router.Router, prefix string) []*route.Route {
	var routes []*route.Route
	prefix += r.Prefix

	for _, route := range r.Routes {
		route.Path = prefix + route.Path
		route.Pattern = route.Protocol + " " + route.Path
		routes = append(routes, &route)
	}

	for _, subRouter := range r.Routers {
		routes = append(routes, a.GetRoutes(subRouter, prefix)...)
	}
	return routes
}

func (a *App) IncludeRouter(r *router.Router) {
	a.RootRouter.AddRouter(r)
}

func (a *App) ListenAndServe() {
	mux := http.NewServeMux()
	router := middleware.LoggingMiddleware(mux)

	routes := a.GetRoutes(a.RootRouter, "")

	for _, route := range routes {
		slog.Info(fmt.Sprintf("Serving route: %s", route.Pattern))
		mux.HandleFunc(route.Pattern, func(w http.ResponseWriter, req *http.Request) {
			handler.Handler(w, req, route)
		})
	}

	var addr string
	if a.Network {
		addr += "0.0.0.0"
	}
	addr += fmt.Sprintf(":%d", a.Port)

	slog.Info(fmt.Sprintf("Running Puff 💨 on port %d", a.Port))

	http.ListenAndServe(addr, router)
}
