package puff

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type Config struct {
	// Name is the application name
	Name string
	// Version is the application version.
	Version string
	// DocsURL is the Router prefix for Swagger documentation. Can be "" to disable Swagger documentation.
	DocsURL string
	// DocsReload, if true, enables automatic reload on the Swagger documentation page.
	DocsReload bool
	// TLSPublicKeyFile specifies the file for the TLS public key (usually .pem or .crt).
	TLSPublicKeyFile string
	// TLSPrivateKeyFile specifies the file for the TLS private key (usually .key).
	TLSPrivateKeyFile string
}

type PuffApp struct {
	// Config is the Puff App Config.
	*Config
	// RootRouter is the application's default router. All routers extend from one.
	RootRouter *Router
	Logger     *slog.Logger
}

// SetDebug sets the application mode to 'DEBUG'.
//
// In this mode, the application will use 'pretty' logging.
func (a *PuffApp) SetDev() {
	logger := a.Logger.Handler().(*SlogHandler)
	logger.SetLevel(slog.LevelDebug)
}

// SetProd sets the application mode to 'PROD'.
//
// In this mode, the application will use structured logging.
func (a *PuffApp) SetProd() {
	handler := a.Logger.Handler().(*SlogHandler)
	handler.SetLevel(slog.LevelInfo)
}

// SetVersion sets the version of the application.
//
// This can be useful for tracking and managing application versions.
func (a *PuffApp) SetVersion(v string) {
	a.Config.Version = v
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

func (a *PuffApp) patchAllRoutes() {
	a.RootRouter.patchRoutes()
	for _, r := range a.RootRouter.Routers {
		r.patchRoutes()
	}
}

func (a *PuffApp) ListenAndServe(listenAddr string) {
	a.patchAllRoutes()
	a.addOpenAPIRoutes()
	slog.Debug(fmt.Sprintf("Running Puff 💨 on %s", listenAddr))
	var err error
	if a.TLSPublicKeyFile != "" && a.TLSPrivateKeyFile != "" {
		err = http.ListenAndServeTLS(listenAddr, a.TLSPublicKeyFile, a.TLSPrivateKeyFile, a.RootRouter)
	} else {
		err = http.ListenAndServe(listenAddr, a.RootRouter)
	}
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func (a *PuffApp) Get(path string, description string, fields any, handleFunc func(*Context)) {
	a.RootRouter.Get(path, description, fields, handleFunc)
}

func (a *PuffApp) Post(path string, description string, fields any, handleFunc func(*Context)) {
	a.RootRouter.Post(path, description, fields, handleFunc)
}

func (a *PuffApp) Patch(path string, description string, fields any, handleFunc func(*Context)) {
	a.RootRouter.Patch(path, description, fields, handleFunc)
}

func (a *PuffApp) Put(path string, description string, fields any, handleFunc func(*Context)) {
	a.RootRouter.Put(path, description, fields, handleFunc)
}

func (a *PuffApp) Delete(path string, description string, fields any, handleFunc func(*Context)) {
	a.RootRouter.Delete(path, description, fields, handleFunc)
}
func (a *PuffApp) WebSocket(path string, description string, fields any, handleFunc func(*Context)) {
	a.RootRouter.WebSocket(path, description, fields, handleFunc)
}

func (a *PuffApp) AllRoutes() []*Route {
	return a.RootRouter.AllRoutes()
}
