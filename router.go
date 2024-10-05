package puff

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

// Router defines a group of routes that share the same prefix and middlewares.
type Router struct {
	Name        string
	Prefix      string //(optional) prefix, all Routes underneath will have paths that start with the prefix automatically
	Routers     []*Router
	Routes      []*Route
	Middlewares []*Middleware
	Tag         string
	Description string
	// Responses is a map of status code to puff.Response. Possible Responses for routes can be set at the Router (root as well),
	// and Route level, however responses directly set on the route will have the highest specificity.
	Responses Responses

	// parent maps to the router's immediate parent. Will be nil for RootRouter
	parent *Router
	// puff maps to the original PuffApp
	puff *PuffApp
}

// NewRouter creates a new router provided router name and path prefix.
func NewRouter(name string, prefix string) *Router {
	return &Router{
		Name:      name,
		Prefix:    prefix,
		Responses: Responses{},
	}
}

func (r *Router) registerRoute(
	description string,
	method string,
	path string,
	handleFunc func(*Context),
	fields any,
) *Route {
	newRoute := Route{
		Description: description,
		Path:        path,
		Handler:     handleFunc,
		Protocol:    method,
		Fields:      fields,
		Router:      r,
		Responses:   Responses{},
	}

	r.Routes = append(r.Routes, &newRoute)
	return &newRoute
}

func (r *Router) Get(
	path string, description string,
	fields any,
	handleFunc func(*Context),
) *Route {
	return r.registerRoute(description, http.MethodGet, path, handleFunc, fields)
}

func (r *Router) Post(
	path string, description string,
	fields any,
	handleFunc func(*Context),
) *Route {
	return r.registerRoute(description, http.MethodPost, path, handleFunc, fields)
}

func (r *Router) Put(
	path string, description string,
	fields any,
	handleFunc func(*Context),
) *Route {
	return r.registerRoute(description, http.MethodPut, path, handleFunc, fields)
}

func (r *Router) Patch(
	path string, description string,
	fields any,
	handleFunc func(*Context),
) *Route {
	return r.registerRoute(description, http.MethodPatch, path, handleFunc, fields)
}

func (r *Router) Delete(
	path string, description string,
	fields any,
	handleFunc func(*Context),
) *Route {
	return r.registerRoute(description, http.MethodDelete, path, handleFunc, fields)
}

func (r *Router) WebSocket(
	path string, description string,
	fields any,
	handleFunc func(*Context),
) *Route {
	newRoute := Route{
		WebSocket: true,
		Protocol:  "GET",
		Path:      path,
		Handler:   handleFunc,
		Fields:    fields,
	}
	r.Routes = append(r.Routes, &newRoute)
	return &newRoute
}

func (r *Router) IncludeRouter(rt *Router) {
	if rt.parent != nil {
		err := fmt.Errorf(
			"provided router is already attached to %s. A router may only be attached to one parent",
			rt.parent,
		)
		panic(err)
	}

	rt.parent = r
	if rt.parent != nil {
		rt.puff = rt.parent.puff
	}
	r.Routers = append(r.Routers, rt)
}

// A
// Root Router -> parent nil, puff a
//
//	PizzaRouter ->  parent RootRouter, puff a
func (r *Router) Use(m Middleware) {
	r.Middlewares = append(r.Middlewares, &m)
}

func (r *Router) String() string {
	return fmt.Sprintf("Name: %s Prefix: %s", r.Name, r.Prefix)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, router := range r.Routers {
		if strings.HasPrefix(req.URL.Path, router.Prefix) {
			router.ServeHTTP(w, req)
			return
		}
	}
	c := NewContext(w, req)
	for _, route := range r.Routes {
		if route.regexp == nil {
			// TODO: need to fix this. this will be nil for the doc routes.
			route.getCompletePath()
			route.createRegexMatch()
		}
		isMatch := route.regexp.MatchString(req.URL.Path)
		if isMatch && req.Method == route.Protocol {
			matches := route.regexp.FindStringSubmatch(req.URL.Path)
			err := populateInputSchema(c, route.Fields, route.params, matches)
			if err != nil {
				c.BadRequest(err.Error())
				return
			}
			if route.WebSocket {
				if !c.isWebSocket() {
					c.BadRequest("This route uses the WebSocket protocol.")
					return
				}
				handleWebSocket(c)
				go c.WebSocket.read()
				handler := route.Handler
				handler(c)
				for c.WebSocket.IsOpen() {
				}
				return
			}
			handler := route.Handler
			handler(c)
			return
		}
	}
	http.NotFound(w, req)
}

func Unprocessable(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "StatusUnprocessableEntity", http.StatusUnprocessableEntity)
}

// AllRoutes returns all routes attached to a router as well as routes attached to the subrouters
// For just the routes attached to a router, use `Routes` attribute on Router
func (r *Router) AllRoutes() []*Route {
	var routes []*Route

	routes = append(routes, r.Routes...)

	for _, subRouter := range r.Routers {
		routes = append(routes, subRouter.AllRoutes()...)
	}
	return routes
}

func (r *Router) patchRoutes() {
	for _, route := range r.Routes {
		route.getCompletePath()
		route.createRegexMatch()
		err := route.handleInputSchema()
		if err != nil {
			panic("error with Input Schema for route " + route.Path + " on router " + r.Name + ". Error: " + err.Error())
		}
		slog.Debug(fmt.Sprintf("Serving route: %s", route.fullPath))
		// populate route with their respective responses
		route.GenerateResponses()
	}
	//TODO: ensure no route collision, will be a nice to have
}
