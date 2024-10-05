package puff

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"reflect"
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
	// OpenAPI configuration. Gives users access to the OpenAPI spec generated. Can be manipulated by the user.
	OpenAPI OpenAPI
}

// Add a Router to the main app.
// Under the hood attaches the router to the App's RootRouter
func (a *PuffApp) IncludeRouter(r *Router) {
	r.puff = a
	a.RootRouter.IncludeRouter(r)
}

func (a *PuffApp) Use(m Middleware) {
	a.RootRouter.Middlewares = append(a.RootRouter.Middlewares, &m)
}

func (a *PuffApp) addOpenAPIRoutes() {
	if a.DocsURL == "" {
		return
	}
	a.GenerateOpenAPISpec()
	docsRouter := Router{
		Prefix: a.DocsURL,
		Name:   "OpenAPI Documentation Router",
	}

	docsRouter.Get(".json", "Provides JSON OpenAPI Schema.", nil, func(c *Context) {
		res := GenericResponse{
			Content:     string(*a.OpenAPI.spec),
			ContentType: "application/json",
		}
		c.SendResponse(res)
	})

	docsRouter.Get("", "Render OpenAPI schema.", nil, func(c *Context) {
		res := HTMLResponse{
			Content: GenerateOpenAPIUI("OpenAPI Spec", a.DocsURL+".json"),
		}
		c.SendResponse(res)
	})
	if a.DocsReload {
		docsRouter.WebSocket("/ws", "WebSocket for live reload of swagger page.", nil, func(c *Context) {
			c.WebSocket.OnMessage = func(ws *WebSocket, wsm WebSocketMessage) {
				msg := new(string)
				err := wsm.To(msg)
				if err != nil {
					ws.Send(err.Error()) // do not care about errs here
					ws.Close()
					return
				}
				if *msg == "ping" {
					ws.Send("pong")
					// if err != nil {
					// 	slog.Debug("pingpong swagger ws: " + err.Error())
					// 	ws.Close()
					// }
				}
				if *msg == "disconnect" {
					ws.Close()
				}
			}
		})
	}
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

func (a *PuffApp) Get(path string, description string, fields any, handleFunc func(*Context)) *Route {
	return a.RootRouter.Get(path, description, fields, handleFunc)
}

func (a *PuffApp) Post(path string, description string, fields any, handleFunc func(*Context)) *Route {
	return a.RootRouter.Post(path, description, fields, handleFunc)
}

func (a *PuffApp) Patch(path string, description string, fields any, handleFunc func(*Context)) *Route {
	return a.RootRouter.Patch(path, description, fields, handleFunc)
}

func (a *PuffApp) Put(path string, description string, fields any, handleFunc func(*Context)) *Route {
	return a.RootRouter.Put(path, description, fields, handleFunc)
}

func (a *PuffApp) Delete(path string, description string, fields any, handleFunc func(*Context)) *Route {
	return a.RootRouter.Delete(path, description, fields, handleFunc)
}
func (a *PuffApp) WebSocket(path string, description string, fields any, handleFunc func(*Context)) *Route {
	return a.RootRouter.WebSocket(path, description, fields, handleFunc)
}

func (a *PuffApp) AllRoutes() []*Route {
	return a.RootRouter.AllRoutes()
}

func (a *PuffApp) SetResponses(r Responses) {
	a.RootRouter.Responses = r
}

func (a *PuffApp) GenerateOpenAPISpec() {
	if reflect.ValueOf(a.OpenAPI).IsZero() {
		paths, tags := a.GeneratePathsTags()
		a.OpenAPI = OpenAPI{
			SpecVersion: "3.1.0",
			Info: Info{
				Version:     a.Version,
				Title:       a.Name,
				Description: "<h4>Application built via Puff Framework</h4>",
			},
			Servers:     []Server{},
			Tags:        tags,
			Paths:       paths,
			Definitions: Definitions,
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
