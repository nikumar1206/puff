package puff

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/nikumar1206/puff/middleware"
)

type Config struct {
	Network bool   // host to the entire network?
	Port    int    // port number to use
	Name    string // title for OpenAPI spec
	Version string // ex. 1.0.0, default: 1.0.0
	DocsURL string
}

type PuffApp struct {
	*Config
	RootRouter  *Router // This is the root router. All other routers will work underneath this.
	Middlewares []*middleware.Middleware
}

// gets all routes for a router
func (a *PuffApp) getRoutes(r *Router, prefix string) []*Route {
	var routes []*Route
	prefix += r.Prefix

	for _, route := range r.Routes {
		route.Path = prefix + route.Path
		route.Pattern = route.Protocol + " " + route.Path
		routes = append(routes, &route)
	}

	for _, subRouter := range r.Routers {
		routes = append(routes, a.getRoutes(subRouter, prefix)...)
	}
	return routes
}

// Add a Router to the main app.
// Under the hood attaches the router to the App's RootRouter
func (a *PuffApp) IncludeRouter(r *Router) {
	a.RootRouter.IncludeRouter(r)
}

func (a *PuffApp) IncludeMiddleware(m middleware.Middleware) {
	a.Middlewares = append(a.Middlewares, &m)
}

func (a *PuffApp) IncludeMiddlewares(ms ...middleware.Middleware) {
	for _, m := range ms {
		a.IncludeMiddleware(m)
	}
}

func (a *PuffApp) AddOpenAPIDocs(mux *http.ServeMux, routes []*Route) {
	if a.DocsURL == "" {
		return
	}
	spec, err := GenerateOpenAPISpec(a.Name, a.Version, routes)
	if err != nil {
		slog.Error(fmt.Sprintf("Generating the OpenAPISpec failed. Error: %s", err.Error()))
		return
	}
	openAPIDocsRoute := Route{
		Protocol:    "GET",
		Path:        a.DocsURL + ".json",
		Pattern:     "GET " + a.DocsURL + ".json",
		Description: "Recieve Docs as JSON.",
		Handler: func(req Request) interface{} {
			res := Response{
				Content: spec,
			}
			res.ContentType = "application/json"
			return res
		},
	}
	openAPIUIDocsRoute := Route{
		Protocol:    "GET",
		Path:        a.DocsURL,
		Pattern:     "GET " + a.DocsURL,
		Description: "Display the OpenAPI Docs in Spotlight.",
		Handler: func(req Request) interface{} {
			return HTMLResponse{
				Content: GenerateOpenAPIUI(spec, "OpenAPI Spec", a.DocsURL+".json"),
			}
		},
	}
	muxAddHandleFunc(mux, &openAPIDocsRoute)
	muxAddHandleFunc(mux, &openAPIUIDocsRoute)
}

// Adds a Route to mux
func muxAddHandleFunc(mux *http.ServeMux, route *Route) {
	mux.HandleFunc(route.Pattern, func(w http.ResponseWriter, req *http.Request) {
		Handler(w, req, route)
	})
}

func (a *PuffApp) ListenAndServe() {
	mux := http.NewServeMux()
	var router http.Handler = mux

	routes := a.getRoutes(a.RootRouter, "")

	for _, route := range routes {
		slog.Info(fmt.Sprintf("Serving route: %s", route.Pattern))
		muxAddHandleFunc(mux, route)
	}

	for _, m := range a.Middlewares {
		router = (*m)(router)
	}

	// Add OpenAPISpec
	a.AddOpenAPIDocs(mux, routes)

	// Listen and Serve
	var addr string
	if a.Network {
		addr += "0.0.0.0"
	}
	addr += fmt.Sprintf(":%d", a.Port)

	slog.Info(fmt.Sprintf("Running Puff 💨 on port %d", a.Port))
	http.ListenAndServe(addr, router)
}

func (a *PuffApp) Get(path string, description string, handleFunc func(Request) interface{}) {
	a.RootRouter.Get(path, description, handleFunc)
}

func (a *PuffApp) Post(path string, description string, handleFunc func(Request) interface{}) {
	a.RootRouter.Post(path, description, handleFunc)
}

func (a *PuffApp) Patch(path string, description string, handleFunc func(Request) interface{}) {
	a.RootRouter.Patch(path, description, handleFunc)
}

func (a *PuffApp) Put(path string, description string, handleFunc func(Request) interface{}) {
	a.RootRouter.Put(path, description, handleFunc)
}

func (a *PuffApp) Delete(path string, description string, handleFunc func(Request) interface{}) {
	a.RootRouter.Delete(path, description, handleFunc)
}

// Gets all Routes for the App
func (a *PuffApp) GetRoutes() []*Route {
	return a.getRoutes(a.RootRouter, "")
}
