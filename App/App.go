package app

import (
	"fmt"
	"log/slog"
	"net/http"
	router "puff/router"
	"time"
)

type Config struct {
	Network bool // host to the entire network?
	Reload  bool // live reload?
	Port    int  // port number to use
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

func getAllRoutes(rtr *router.Router, prefix string) []string {
	var routes []string
	for _, route := range rtr.Routes {
		routes = append(routes, fmt.Sprintf("%s %s%s", route.Protocol, prefix, route.Path)) //ex: GET /food/cheese/swiss
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

	mux.HandleFunc("/", a.sendToHandler)

	addr := ""
	if a.Network {
		addr += "0.0.0.0"
	}
	addr += fmt.Sprintf(":%d", a.Port)

	slog.Info(fmt.Sprintf("Running Puff 💨 on port %d", a.Port))
	http.ListenAndServe(addr, router)
}
