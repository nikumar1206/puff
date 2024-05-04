package router

import (
	"net/http"

	"github.com/nikumar1206/puff/request"
	"github.com/nikumar1206/puff/route"
)

type Router struct {
	Name    string
	Prefix  string //(optional) prefix, all Routes underneath will have paths that start with the prefix automatically
	Routers []*Router
	Routes  []route.Route
	// middlewares []Middleware
}

func (r *Router) registerRoute(
	method string,
	path string,
	handleFunc func(request.Request) interface{},
) {
	newRoute := route.Route{
		RouterName: r.Name,
		Path:       path,
		Handler:    handleFunc,
		Protocol:   method,
		Pattern:    method + " " + path,
	}
	r.Routes = append(r.Routes, newRoute)
}
func (r *Router) GET(
	path string,
	description string,
	handleFunc func(request.Request) interface{},
) {
	r.registerRoute(http.MethodGet, path, handleFunc)
}

func (r *Router) POST(
	path string,
	description string,
	handleFunc func(request.Request) interface{},
) {
	r.registerRoute(http.MethodPost, path, handleFunc)
}

func (r *Router) PUT(
	path string,
	description string,
	handleFunc func(request.Request) interface{},
) {
	r.registerRoute(http.MethodPut, path, handleFunc)
}

func (r *Router) PATCH(
	path string,
	description string,
	handleFunc func(request.Request) interface{},
) {
	r.registerRoute(http.MethodPatch, path, handleFunc)
}

func (r *Router) IncludeRouter(rt *Router) *Router {
	r.Routers = append(r.Routers, rt)
	return rt
}
