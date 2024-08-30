package puff

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type PuffApp struct {
	// Name is the application name
	Name string
	// Version is the application version.
	Version string
	// DocsURL is the Router prefix for Swagger documentation. Can be "" to disable Swagger documentation.
	DocsURL string
	// DocsReload, if true, enables automatic reload on the Swagger documentation page.
	DocsReload bool
	// TLSPublicCertFile specifies the file for the TLS certificate (usually .pem or .crt).
	TLSPublicCertFile string
	// TLSPrivateKeyFile specifies the file for the TLS private key (usually .key).
	TLSPrivateKeyFile string
	// RootRouter is the application's default router. All routers extend from one.
	RootRouter *Router
	// Logger is the reference to the application's logger. Equivalent to slog.Default()
	Logger *slog.Logger
}

// Add a Router to the main app.
// Under the hood attaches the router to the App's RootRouter
func (a *PuffApp) IncludeRouter(r *Router) {
	a.RootRouter.IncludeRouter(r)
}

// Use registers a middleware function to be used by the root router of the PuffApp.
// The middleware will be appended to the list of middlewares in the root router.
//
// Parameters:
// - m: Middleware function to be added.
func (a *PuffApp) Use(m Middleware) {
	a.RootRouter.Middlewares = append(a.RootRouter.Middlewares, &m)
}

// addOpenAPIRoutes adds routes to serve OpenAPI documentation for the PuffApp.
// If a DocsURL is specified, the function sets up two routes:
// 1. A route to provide the OpenAPI spec as JSON.
// 2. A route to render the OpenAPI documentation in a user-friendly UI.
//
// This method will not add any routes if DocsURL is empty.
//
// Errors during spec generation are logged, and the method will exit early if any occur.
func (a *PuffApp) addOpenAPIRoutes() {
	if a.DocsURL == "" {
		return
	}
	spec, err := GenerateOpenAPISpec(a.Name, a.Version, *a.RootRouter)
	if err != nil {
		slog.Error(fmt.Sprintf("Generating the OpenAPISpec failed. Error: %s", err.Error()))
		return
	}
	docsRouter := Router{
		Prefix: a.DocsURL,
		Name:   "OpenAPI Documentation Router",
	}

	// Provides JSON OpenAPI Schema.
	docsRouter.Get(".json", nil, func(c *Context) {
		res := GenericResponse{
			Content:     spec,
			ContentType: "application/json",
		}
		c.SendResponse(res)
	})

	// Renders OpenAPI schema.
	docsRouter.Get("", nil, func(c *Context) {
		res := HTMLResponse{
			Content: GenerateOpenAPIUI(spec, "OpenAPI Spec", a.DocsURL+".json"),
		}
		c.SendResponse(res)
	})

	a.IncludeRouter(&docsRouter)
}

// attachMiddlewares recursively applies middlewares to all routes within a router.
// This function traverses through the router's sub-routers and routes, applying the
// middleware functions in the given order.
//
// Parameters:
// - middleware_combo: A pointer to a slice of Middleware to be applied.
// - router: The router whose middlewares and routes should be processed.
func attachMiddlewares(middleware_combo *[]Middleware, router *Router) {
	for _, m := range router.Middlewares {
		nmc := append(*middleware_combo, *m)
		middleware_combo = &nmc
	}
	for _, route := range router.Routes {
		for _, m := range *middleware_combo {
			route.Handler = (m)(route.Handler)
		}
	}
	for _, router := range router.Routers {
		attachMiddlewares((middleware_combo), router)
	}
}

// patchAllRoutes applies middlewares to all routes and sub-routers in the root router
// of the PuffApp. It also patches the routes of each router to ensure they have been
// processed for middlewares.
func (a *PuffApp) patchAllRoutes() {
	a.RootRouter.patchRoutes()
	for _, r := range a.RootRouter.Routers {
		r.patchRoutes()
	}
	attachMiddlewares(&[]Middleware{}, a.RootRouter)
}

// ListenAndServe starts the PuffApp server on the specified address.
// Before starting, it patches all routes, adds OpenAPI documentation routes (if available),
// and sets up logging.
//
// If TLS certificates are provided (TLSPublicCertFile and TLSPrivateKeyFile), the server
// starts with TLS enabled; otherwise, it runs a standard HTTP server.
//
// Parameters:
// - listenAddr: The address the server will listen on (e.g., ":8080").
func (a *PuffApp) ListenAndServe(listenAddr string) {
	slog.SetDefault(a.Logger)
	a.patchAllRoutes()
	a.addOpenAPIRoutes()
	slog.Debug(fmt.Sprintf("Running Puff 💨 on %s", listenAddr))
	var err error
	if a.TLSPublicCertFile != "" && a.TLSPrivateKeyFile != "" {
		err = http.ListenAndServeTLS(listenAddr, a.TLSPublicCertFile, a.TLSPrivateKeyFile, a.RootRouter)
	} else {
		err = http.ListenAndServe(listenAddr, a.RootRouter)
	}
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

// Get registers an HTTP GET route in the PuffApp's root router.
//
// Parameters:
// - path: The URL path of the route.
// - fields: Optional fields associated with the route.
// - handleFunc: The handler function that will be executed when the route is accessed.
func (a *PuffApp) Get(path string, fields any, handleFunc func(*Context)) {
	a.RootRouter.registerRoute(http.MethodGet, path, handleFunc, fields)
}

// Post registers an HTTP POST route in the PuffApp's root router.
//
// Parameters:
// - path: The URL path of the route.
// - fields: Optional fields associated with the route.
// - handleFunc: The handler function that will be executed when the route is accessed.
func (a *PuffApp) Post(path string, fields any, handleFunc func(*Context)) {
	a.RootRouter.registerRoute(http.MethodPost, path, handleFunc, fields)
}

// Patch registers an HTTP PATCH route in the PuffApp's root router.
//
// Parameters:
// - path: The URL path of the route.
// - fields: Optional fields associated with the route.
// - handleFunc: The handler function that will be executed when the route is accessed.
func (a *PuffApp) Patch(path string, fields any, handleFunc func(*Context)) {
	a.RootRouter.registerRoute(http.MethodPatch, path, handleFunc, fields)
}

// Put registers an HTTP PUT route in the PuffApp's root router.
//
// Parameters:
// - path: The URL path of the route.
// - fields: Optional fields associated with the route.
// - handleFunc: The handler function that will be executed when the route is accessed.
func (a *PuffApp) Put(path string, fields any, handleFunc func(*Context)) {
	a.RootRouter.registerRoute(http.MethodPut, path, handleFunc, fields)
}

// Delete registers an HTTP DELETE route in the PuffApp's root router.
//
// Parameters:
// - path: The URL path of the route.
// - fields: Optional fields associated with the route.
// - handleFunc: The handler function that will be executed when the route is accessed.
func (a *PuffApp) Delete(path string, fields any, handleFunc func(*Context)) {
	a.RootRouter.registerRoute(http.MethodDelete, path, handleFunc, fields)
}

// WebSocket registers a WebSocket route in the PuffApp's root router.
// This route allows the server to handle WebSocket connections at the specified path.
//
// Parameters:
// - path: The URL path of the WebSocket route.
// - fields: Optional fields associated with the route.
// - handleFunc: The handler function to handle WebSocket connections.
func (a *PuffApp) WebSocket(path string, fields any, handleFunc func(*Context)) {
	newRoute := Route{
		WebSocket: true,
		Protocol:  "GET",
		Path:      path,
		Handler:   handleFunc,
		Fields:    fields,
	}
	a.RootRouter.Routes = append(a.RootRouter.Routes, &newRoute)
}

// AllRoutes returns all routes registered in the PuffApp, including those in sub-routers.
// This function provides an aggregated view of all routes in the application.
func (a *PuffApp) AllRoutes() []*Route {
	return a.RootRouter.AllRoutes()
}
