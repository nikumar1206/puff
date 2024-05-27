package puff

import (
	"fmt"
	"net/http"
)

type Router struct {
	Name    string
	Prefix  string //(optional) prefix, all Routes underneath will have paths that start with the prefix automatically
	Routers []*Router
	Routes  []Route
	// middlewares []Middleware
}

func (r *Router) registerRoute(
	method string,
	path string,
	handleFunc func(Request) Response,
	description string,
) {
	newRoute := Route{
		RouterName:  r.Name,
		Description: description,
		Path:        path,
		Handler:     handleFunc,
		Protocol:    method,
		Pattern:     method + " " + path,
	}
	r.Routes = append(r.Routes, newRoute)
	r.Routes = append(r.Routes, newRoute)
}

func (r *Router) Get(
	path string,
	description string,
	handleFunc func(Request) Response,
) {
	r.registerRoute(http.MethodGet, path, handleFunc, description)
}

func (r *Router) Post(
	path string,
	description string,
	handleFunc func(Request) Response,
) {
	r.registerRoute(http.MethodPost, path, handleFunc, description)
}

func (r *Router) Put(
	path string,
	description string,
	handleFunc func(Request) Response,
) {
	r.registerRoute(http.MethodPut, path, handleFunc, description)
}

func (r *Router) Patch(
	path string,
	description string,
	handleFunc func(Request) Response,
) {
	r.registerRoute(http.MethodPatch, path, handleFunc, description)
}

func (r *Router) Delete(
	path string,
	description string,
	handleFunc func(Request) Response,
) {
	r.registerRoute(http.MethodDelete, path, handleFunc, description)
}

func (r *Router) IncludeRouter(rt *Router) {
	r.Routers = append(r.Routers, rt)
}

func (r *Router) String() string {
	return fmt.Sprintf("Name: %s Prefix: %s", r.Name, r.Prefix)
}
