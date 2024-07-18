package puff

import (
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
)

// Router defines a group of routes that share the same prefix and middlewares.
type Router struct {
	Name        string
	Prefix      string //(optional) prefix, all Routes underneath will have paths that start with the prefix automatically
	Routers     []*Router
	Routes      []*Route
	Middlewares []*Middleware
	parent      *Router
	Tag         string
	Description string
}

// NewRouter creates a new router provided router name and path prefix.
func NewRouter(name string, prefix string) *Router {
	return &Router{
		Name:   name,
		Prefix: prefix,
	}
}

func (r *Router) registerRoute(
	method string,
	path string,
	handleFunc func(*Context),
	fields any,
) {
	newRoute := Route{
		Path:     path,
		Handler:  handleFunc,
		Protocol: method,
		Fields:   fields,
	}

	r.Routes = append(r.Routes, &newRoute)
}

func (r *Router) Get(
	path string,
	fields any,
	handleFunc func(*Context),
) {
	r.registerRoute(http.MethodGet, path, handleFunc, fields)
}

func (r *Router) Post(
	path string,
	fields any,
	handleFunc func(*Context),
) {
	r.registerRoute(http.MethodPost, path, handleFunc, fields)
}

func (r *Router) Put(
	path string,
	fields any,
	handleFunc func(*Context),
) {
	r.registerRoute(http.MethodPut, path, handleFunc, fields)
}

func (r *Router) Patch(
	path string,
	fields any,
	handleFunc func(*Context),
) {
	r.registerRoute(http.MethodPatch, path, handleFunc, fields)
}

func (r *Router) Delete(
	path string,
	fields any,
	handleFunc func(*Context),
) {
	r.registerRoute(http.MethodDelete, path, handleFunc, fields)
}

func (r *Router) WebSocket(
	path string,
	fields any,
	handleFunc func(*Context),
) {
	newRoute := Route{
		WebSocket: true,
		Protocol:  "GET",
		Path:      path,
		Handler:   handleFunc,
		Fields:    fields,
	}
	r.Routes = append(r.Routes, &newRoute)
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
	r.Routers = append(r.Routers, rt)
}

func (r *Router) String() string {
	return fmt.Sprintf("Name: %s Prefix: %s", r.Name, r.Prefix)
}

func (r *Router) getCompletePath(route *Route) {
	var parts []string
	currentRouter := r
	for currentRouter != nil {
		parts = append([]string{currentRouter.Prefix}, parts...)
		currentRouter = currentRouter.parent
	}

	parts = append(parts, route.Path)
	route.fullPath = strings.Join(parts, "")
}

func (r *Router) createRegexMatch(route *Route) {
	escapedPath := strings.ReplaceAll(route.fullPath, "/", "\\/")
	regexPattern := regexp.MustCompile(`\{[^}]+\}`).ReplaceAllString(escapedPath, "([^/]+)")
	route.regexp = regexp.MustCompile("^" + regexPattern + "$")
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
			r.getCompletePath(route)
			r.createRegexMatch(route)
		}
		isMatch := route.regexp.MatchString(req.URL.Path)
		if isMatch && req.Method == route.Protocol {
			// err := route.anys.ValidateIncomingAttribute(any.Responses, "cheese")
			// if err != nil {
			// 	Unprocessable(w, req)
			// }
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
				for _, m := range r.Middlewares {
					handler = (*m)(handler)
				}
				handler(c)
				for c.WebSocket.IsOpen() {
				}
			}
			handler := route.Handler
			for _, m := range r.Middlewares {
				handler = (*m)(handler)
			}
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
		r.getCompletePath(route)
		r.createRegexMatch(route)
		err := handleInputSchema(&route.params, route.Fields)
		if err != nil {
			panic("Error with Input Schema for route " + route.Path + " on router " + r.Name + ". Error: " + err.Error())
		}
		slog.Debug(fmt.Sprintf("Serving route: %s", route.fullPath))
	}
	//TODO: ensure no route collision, will be a nice to have
}
