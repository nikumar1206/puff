package puff

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type Config struct {
	// ListenAddr is the address to listen on.
	// TODO: Remove ListenAddr, specify as argument to listen and serve.
	ListenAddr string
	// Name is the application name
	Name string
	// Version is the application version.
	Version string
	// DocsURL is the Router prefix for Swagger documentation. Can be "" to disable Swagger documentation.
	DocsURL string
	// TODO: depending on the mode, set the log level and other settings
	Mode string
}

type PuffApp struct {
	*Config
	RootRouter  *Router // This is the root router. All other routers will work underneath this.
	Middlewares []*Middleware
}

// Add a Router to the main app.
// Under the hood attaches the router to the App's RootRouter

func (a *PuffApp) IncludeRouter(r *Router) {
	a.RootRouter.IncludeRouter(r)
}

func (a *PuffApp) IncludeMiddleware(m Middleware) {
	a.Middlewares = append(a.Middlewares, &m)
}

func (a *PuffApp) IncludeMiddlewares(ms ...Middleware) {
	for _, m := range ms {
		a.IncludeMiddleware(m)
	}
}

func (a *PuffApp) AddOpenAPIRoutes() {
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

	docsRouter.Get(".json", Field{}, func(c *Context) {
		res := GenericResponse{
			Content:     spec,
			ContentType: "application/json",
		}
		c.SendResponse(res)
	})

	docsRouter.Get("", Field{}, func(c *Context) {
		res := HTMLResponse{
			Content: GenerateOpenAPIUI(spec, "OpenAPI Spec", a.DocsURL+".json"),
		}
		c.SendResponse(res)
	})

	a.IncludeRouter(&docsRouter)
}

func (a *PuffApp) ListenAndServe() {

	a.AddOpenAPIRoutes()
	a.RootRouter.patchRoutes()

	for _, r := range a.RootRouter.Routers {
		r.patchRoutes()

		for _, route := range r.Routes {
			slog.Info(fmt.Sprintf("Serving route: %s", route.fullPath))
		}
	}

	slog.Info(fmt.Sprintf("Running Puff 💨 on %s", a.ListenAddr))

	err := http.ListenAndServe(a.ListenAddr, a.RootRouter)

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func (a *PuffApp) Get(path string, fields Field, handleFunc func(*Context)) {
	a.RootRouter.Get(path, fields, handleFunc)
}

func (a *PuffApp) Post(path string, fields Field, handleFunc func(*Context)) {
	a.RootRouter.Post(path, fields, handleFunc)
}

func (a *PuffApp) Patch(path string, fields Field, handleFunc func(*Context)) {
	a.RootRouter.Patch(path, fields, handleFunc)
}

func (a *PuffApp) Put(path string, fields Field, handleFunc func(*Context)) {
	a.RootRouter.Put(path, fields, handleFunc)
}

func (a *PuffApp) Delete(path string, fields Field, handleFunc func(*Context)) {
	a.RootRouter.Delete(path, fields, handleFunc)
}
func (a *PuffApp) WebSocket(path string, fields Field, handleFunc func(*Context)) {
	a.RootRouter.WebSocket(path, fields, handleFunc)
}

func (a *PuffApp) AllRoutes() []*Route {
	return a.RootRouter.AllRoutes()
}
