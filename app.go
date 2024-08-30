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

func (a *PuffApp) Use(m Middleware) {
	a.RootRouter.Middlewares = append(a.RootRouter.Middlewares, &m)
}

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

	docsRouter.Get(".json", "Provides JSON OpenAPI Schema.", nil, func(c *Context) {
		res := GenericResponse{
			Content:     spec,
			ContentType: "application/json",
		}
		c.SendResponse(res)
	})

	docsRouter.Get("", "Render OpenAPI schema.", nil, func(c *Context) {
		res := HTMLResponse{
			Content: GenerateOpenAPIUI(spec, "OpenAPI Spec", a.DocsURL+".json"),
		}
		c.SendResponse(res)
	})

	a.IncludeRouter(&docsRouter)
}

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

func (a *PuffApp) patchAllRoutes() {
	a.RootRouter.patchRoutes()
	for _, r := range a.RootRouter.Routers {
		r.patchRoutes()
	}
	attachMiddlewares(&[]Middleware{}, a.RootRouter)
}

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

func (a *PuffApp) Get(path string, description string, fields any, handleFunc func(*Context)) {
	a.RootRouter.registerRoute(description, http.MethodGet, path, handleFunc, fields)
}

func (a *PuffApp) Post(path string, description string, fields any, handleFunc func(*Context)) {
	a.RootRouter.registerRoute(description, http.MethodPost, path, handleFunc, fields)
}

func (a *PuffApp) Patch(path string, description string, fields any, handleFunc func(*Context)) {
	a.RootRouter.registerRoute(description, http.MethodPatch, path, handleFunc, fields)
}

func (a *PuffApp) Put(path string, description string, fields any, handleFunc func(*Context)) {
	a.RootRouter.registerRoute(description, http.MethodPut, path, handleFunc, fields)
}

func (a *PuffApp) Delete(path string, description string, fields any, handleFunc func(*Context)) {
	a.RootRouter.registerRoute(description, http.MethodDelete, path, handleFunc, fields)
}
func (a *PuffApp) WebSocket(path string, description string, fields any, handleFunc func(*Context)) {
	newRoute := Route{
		WebSocket: true,
		Protocol:  "GET",
		Path:      path,
		Handler:   handleFunc,
		Fields:    fields,
	}
	a.RootRouter.Routes = append(a.RootRouter.Routes, &newRoute)
}

func (a *PuffApp) AllRoutes() []*Route {
	return a.RootRouter.AllRoutes()
}
