package puff

import (
	"fmt"
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
	handleFunc func(*Context, *interface{}),
	description string,
) {
	newRoute := Route{
		RouterName:  r.Name,
		Description: description,
		Path:        path,
		Handler:     handleFunc,
		Protocol:    method,
		Pattern:     fmt.Sprintf("%s %s", method, path),
	}
	r.Routes = append(r.Routes, newRoute)
}

func (r *Router) Get(path, description string, handleFunc func(*Context, *interface{})) {
	r.registerRoute("GET", path, handleFunc, description)
}
func (r *Router) Post(path, description string, handleFunc func(*Context, *interface{})) {
	r.registerRoute("POST", path, handleFunc, description)
}
func (r *Router) Put(path, description string, handleFunc func(*Context, *interface{})) {
	r.registerRoute("PUT", path, handleFunc, description)
}
func (r *Router) Patch(path, description string, handleFunc func(*Context, *interface{})) {
	r.registerRoute("PATCH", path, handleFunc, description)
}
func (r *Router) Delete(path, description string, handleFunc func(*Context, *interface{})) {
	r.registerRoute("DELETE", path, handleFunc, description)
}

func (r *Router) IncludeRouter(rt *Router) {
	r.Routers = append(r.Routers, rt)
}

func (r *Router) String() string {
	return fmt.Sprintf("Name: %s Prefix: %s", r.Name, r.Prefix)
}
