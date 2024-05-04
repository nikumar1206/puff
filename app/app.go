package app

import (
	"fmt"
	"log/slog"
	"net/http"
	handler "puff/handler"
	openapi "puff/openapi"
	"puff/request"
	"puff/response"
	route "puff/route"
	router "puff/router"
	"time"
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
	// Routes  []route.Route
	Router router.Router //This is the root router. All other routers will work underneath this.
	// Middlewares
}

func (a *App) IncludeRouter(r *router.Router) {
	a.Router.IncludeRouter(r)
}
func getAllRoutes(rtr *router.Router, prefix string) []*route.Route {
	var routes []*route.Route
	for _, route := range rtr.Routes {
		route.Path = fmt.Sprintf("%s%s", prefix, route.Path)
		route.Pattern = fmt.Sprintf("%s %s", route.Protocol, route.Path) //ex: GET /food/cheese/swiss
		if route.Protocol == "GET" && route.Fields != nil {
			for _, field := range route.Fields {
				route.Pattern += fmt.Sprintf("/{%s}", field.Name)
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
	router := loggingMiddleware(mux)

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
