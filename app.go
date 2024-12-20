package puff

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
)

type PuffApp struct {
	// Name is the application name
	Name string
	// Version is the application version.
	Version string
	// DocsURL is the Router prefix for Swagger documentation. Can be "" to disable Swagger documentation.
	DocsURL string
	// TLSPublicCertFile specifies the file for the TLS certificate (usually .pem or .crt).
	TLSPublicCertFile string
	// TLSPrivateKeyFile specifies the file for the TLS private key (usually .key).
	TLSPrivateKeyFile string
	// RootRouter is the application's default router. All routers extend from one.
	RootRouter *Router
	// Logger is the reference to the application's logger. Equivalent to slog.Default()
	Logger *slog.Logger
	// OpenAPI configuration. Gives users access to the OpenAPI spec generated. Can be manipulated by the user.
	OpenAPI *OpenAPI
	// TLSConfig to pass into the underlying http.Server
	TLSConfig *tls.Config
	// the underlying server that powers Puff.
	server *http.Server
}

// Add a Router to the main app.
// Under the hood attaches the router to the App's RootRouter
func (a *PuffApp) IncludeRouter(r *Router) {
	r.puff = a
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
	a.GenerateOpenAPISpec()
	docsRouter := Router{
		Prefix: a.DocsURL,
		Name:   "OpenAPI Documentation Router",
	}

	// Provides JSON OpenAPI Schema.
	docsRouter.Get(".json", nil, func(c *Context) {
		res := GenericResponse{
			Content:     string(*a.OpenAPI.spec),
			ContentType: "application/json",
		}
		c.SendResponse(res)
	})

	// Renders OpenAPI schema.
	docsRouter.Get("", nil, func(c *Context) {
		if a.OpenAPI.SwaggerUIConfig == nil {

			swaggerConfig := SwaggerUIConfig{
				Title:           a.Name,
				URL:             a.DocsURL + ".json",
				Theme:           "obsidian",
				Filter:          true,
				RequestDuration: false,
				FaviconURL:      "https://fav.farm/💨",
			}
			a.OpenAPI.SwaggerUIConfig = &swaggerConfig
		}
		res := HTMLResponse{
			Template: openAPIHTML, Data: a.OpenAPI.SwaggerUIConfig,
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
func (a *PuffApp) ListenAndServe(listenAddr string) error {
	slog.SetDefault(a.Logger)
	a.patchAllRoutes()
	a.addOpenAPIRoutes()
	slog.Debug(fmt.Sprintf("Running Puff 💨 on %s", listenAddr))
	slog.Debug(fmt.Sprintf("Visit docs 💨 on %s", fmt.Sprintf("http://localhost%s%s", listenAddr, a.DocsURL)))
	var err error

	httpServer := &http.Server{
		Addr:      listenAddr,
		Handler:   a.RootRouter,
		TLSConfig: a.TLSConfig,
	} // TODO: allow setting server level read/write/connect timeouts and middleware/route level.
	a.server = httpServer

	if a.TLSPublicCertFile != "" && a.TLSPrivateKeyFile != "" {
		err = a.server.ListenAndServeTLS(a.TLSPublicCertFile, a.TLSPrivateKeyFile)
	} else {
		err = a.server.ListenAndServe()
	}

	return err
}

// Get registers an HTTP GET route in the PuffApp's root router.
//
// Parameters:
// - path: The URL path of the route.
// - fields: Optional fields associated with the route.
// - handleFunc: The handler function that will be executed when the route is accessed.
func (a *PuffApp) Get(path string, fields any, handleFunc func(*Context)) *Route {
	return a.RootRouter.Get(path, fields, handleFunc)
}

// Post registers an HTTP POST route in the PuffApp's root router.
//
// Parameters:
// - path: The URL path of the route.
// - fields: Optional fields associated with the route.
// - handleFunc: The handler function that will be executed when the route is accessed.
func (a *PuffApp) Post(path string, fields any, handleFunc func(*Context)) *Route {
	return a.RootRouter.Post(path, fields, handleFunc)
}

// Patch registers an HTTP PATCH route in the PuffApp's root router.
//
// Parameters:
// - path: The URL path of the route.
// - fields: Optional fields associated with the route.
// - handleFunc: The handler function that will be executed when the route is accessed.
func (a *PuffApp) Patch(path string, fields any, handleFunc func(*Context)) *Route {
	return a.RootRouter.Patch(path, fields, handleFunc)
}

// Put registers an HTTP PUT route in the PuffApp's root router.
//
// Parameters:
// - path: The URL path of the route.
// - fields: Optional fields associated with the route.
// - handleFunc: The handler function that will be executed when the route is accessed.
func (a *PuffApp) Put(path string, fields any, handleFunc func(*Context)) *Route {
	return a.RootRouter.Put(path, fields, handleFunc)
}

// Delete registers an HTTP DELETE route in the PuffApp's root router.
//
// Parameters:
// - path: The URL path of the route.
// - fields: Optional fields associated with the route.
// - handleFunc: The handler function that will be executed when the route is accessed.
func (a *PuffApp) Delete(path string, fields any, handleFunc func(*Context)) *Route {
	return a.RootRouter.Delete(path, fields, handleFunc)
}

// WebSocket registers a WebSocket route in the PuffApp's root router.
// This route allows the server to handle WebSocket connections at the specified path.
//
// Parameters:
// - path: The URL path of the WebSocket route.
// - fields: Optional fields associated with the route.
// - handleFunc: The handler function to handle WebSocket connections.
func (a *PuffApp) WebSocket(path string, fields any, handleFunc func(*Context)) *Route {
	return a.RootRouter.WebSocket(path, fields, handleFunc)

}

// AllRoutes returns all routes registered in the PuffApp, including those in sub-routers.
// This function provides an aggregated view of all routes in the application.
func (a *PuffApp) AllRoutes() []*Route {
	return a.RootRouter.AllRoutes()
}

func (a *PuffApp) GenerateOpenAPISpec() {
	if reflect.ValueOf(a.OpenAPI).IsZero() {
		paths, tags := a.GeneratePathsTags()
		a.OpenAPI = &OpenAPI{
			SpecVersion: "3.1.0",
			Info: Info{
				Version:     a.Version,
				Title:       a.Name,
				Description: "<h4>Application built via Puff Framework</h4>",
			},
			Servers:  []Server{},
			Tags:     tags,
			Paths:    paths,
			Security: []SecurityRequirement{},
			Webhooks: map[string]any{},
			Components: Components{
				Schemas:         Schemas,
				Responses:       make(map[string]any),
				Parameters:      make(map[string]any),
				Examples:        make(map[string]any),
				RequestBodies:   make(map[string]any),
				SecuritySchemes: make(map[string]any),
				Headers:         make(map[string]any),
				Callbacks:       make(map[string]any),
				PathItems:       make(map[string]any),
				Links:           make(map[string]any),
			},
		}
	}
	// this value is hardcoded. it cannot be changed
	a.OpenAPI.SpecVersion = "3.1.0"
	openAPISpec, err := json.Marshal(a.OpenAPI)
	if err != nil {
		panic(err)
	}
	a.OpenAPI.spec = &openAPISpec
}

// GeneratePathsTags is a helper function to auto-define OpenAPI tags and paths if you would like to customize OpenAPI schema.
// Returns (paths, tagss) to populate the 'Paths' and 'Tags' attribute of OpenAPI
func (a *PuffApp) GeneratePathsTags() (Paths, []Tag) {
	var tags []Tag
	var tagNames []string
	var paths = make(Paths)
	for _, route := range a.RootRouter.Routes {
		addRoute(route, &tags, &tagNames, &paths)
	}
	for _, router := range a.RootRouter.Routers {
		for _, route := range router.Routes {
			addRoute(route, &tags, &tagNames, &paths)
		}
	}
	return paths, tags
}

// GenerateDefinitions is a helper function to auto-define OpenAPI tags and paths if you would like to customize OpenAPI schema.
// Returns (paths, tagss) to populate the 'Paths' and 'Tags' attribute of OpenAPI
func (a *PuffApp) GenerateDefinitions(paths Paths) map[string]*Schema {

	definitions := map[string]*Schema{}
	for _, p := range paths {
		for _, routeParams := range p.Parameters {
			definitions[routeParams.Name] = &routeParams.Schema
		}

	}

	return definitions
}

// Shutdown calls shutdown on the underlying server with a non-nil empty context.
func (a *PuffApp) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

// Close calls close on the underlying server.
func (a *PuffApp) Close() error {
	return a.server.Close()
}
