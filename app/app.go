package app

import (
	"fmt"
	"log/slog"
	"net/http"
	handler "puff/handler"
	openapi "puff/openapi"
	route "puff/route"
	router "puff/router"
	"time"
)

type Config struct {
	Network bool // host to the entire network?
	Port    int  // port number to use
	OpenAPI *openapi.OpenAPI
}

type App struct {
	*Config
	// Routes  []route.Route
	Router router.Router //This is the root router. All other routers will work underneath this.
	// Middlewares
}

func (a *App) IncludeRouter(r *router.Router) {
	a.Router.Add(r)
}

func getAllRoutes(rtr *router.Router, prefix string) []*route.Route {
	var routes []*route.Route
	for _, route := range rtr.Routes {
		route.Path = fmt.Sprintf("%s %s%s", route.Protocol, prefix, route.Path) //ex: GET /food/cheese/swiss
		if route.Protocol == "GET" && route.Fields != nil {
			for _, field := range route.Fields {
				route.Path += fmt.Sprintf("/{%s}", field.Name)
			}
		}
		routes = append(routes, &route)
	}
	for _, rt := range rtr.Routers {
		prefix += rt.Prefix
		routes = append(routes, getAllRoutes(rt, prefix)...)
	}
	return routes
}
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		processingTime := time.Since(startTime).String()
		slog.Info(
			"HTTP Request",
			slog.String("HTTP METHOD", r.Method),
			slog.String("URL", r.URL.String()),
			slog.String("Processing Time", processingTime),
		)
	})
}

func (a *App) ListenAndServe() {
	mux := http.NewServeMux()
	router := loggingMiddleware(mux)

	routes := getAllRoutes(&a.Router, a.Router.Prefix)
	for _, route := range routes {
		mux.HandleFunc(route.Path, func(w http.ResponseWriter, req *http.Request) {
			handler.Handler(w, req, route)
		})
	}
	addr := ""
	if a.Network {
		addr += "0.0.0.0"
	}
	addr += fmt.Sprintf(":%d", a.Port)

	slog.Info(fmt.Sprintf("Running Puff 💨 on port %d", a.Port))

	http.ListenAndServe(addr, router)
}
