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
	Network     bool   // host to the entire network?
	Port        int    // port number to use
	Name        string //title for OpenAPI spec
	Version     string //ex. 1.0.0, default: 1.0.0
	OpenAPIDocs bool
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

func (a *App) AddOpenAPIDocs(mux *http.ServeMux, routes []*route.Route) {
	if a.OpenAPIDocs {
		spec, err := openapi.GenerateOpenAPISpec(a.Name, a.Version, routes)
		fmt.Println(spec)
		if err != nil {
			slog.Error("Generating the OpenAPISpec failed. Error: %s", err.Error())
			return
		}
		openAPIDocsRoute := route.Route{
			Protocol:    "GET",
			Path:        "/api/docs/docs.json",
			Pattern:     "GET /api/docs/docs.json",
			Description: "Recieve Docs as JSON.",
			Fields:      nil,
			Handler: func(req request.Request) interface{} {
				res := response.Response{
					Content: spec,
				}
				res.ContentType = "application/json"
				return res
			},
		}
		openAPIUIDocsRoute := route.Route{
			Protocol:    "GET",
			Path:        "/api/docs",
			Pattern:     "GET /api/docs",
			Description: "Display the OpenAPI Docs in Spotlight.",
			Fields:      nil,
			Handler: func(req request.Request) interface{} {
				return response.HTMLResponse{
					Content: openapi.GenerateOpenAPIUI(spec, "OpenAPI Spec"),
				}
			},
		}
		muxAddHandleFunc(mux, &openAPIDocsRoute)
		muxAddHandleFunc(mux, &openAPIUIDocsRoute)
	}
}

// Adds a route.Route to mux
func muxAddHandleFunc(mux *http.ServeMux, route *route.Route) {
	mux.HandleFunc(route.Pattern, func(w http.ResponseWriter, req *http.Request) {
		handler.Handler(w, req, route)
	})
}

func (a *App) ListenAndServe() {
	mux := http.NewServeMux()
	router := middleware.LoggingMiddleware(mux)

	routes := a.GetRoutes(a.RootRouter, "")

	//Handle Routing
	routes := getAllRoutes(&a.Router, a.Router.Prefix)

	for _, route := range routes {
		muxAddHandleFunc(mux, route)
	}

	//Add OpenAPISpec
	a.AddOpenAPIDocs(mux, routes)

	//Listen and Serve
	addr := ""
	if a.Network {
		addr += "0.0.0.0"
	}
	addr += fmt.Sprintf(":%d", a.Port)

	slog.Info(fmt.Sprintf("Running Puff 💨 on port %d", a.Port))

	http.ListenAndServe(addr, router)
}
