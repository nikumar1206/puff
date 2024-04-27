package router

import (
	field "puff/field"
	route "puff/route"
)

type Router struct {
	prefix  string //(optional) prefix, all routes underneath will have paths that start with the prefix automatically
	routers []*Router
	routes  []route.Route
	// middlewares []Middleware
}

func (a *Router) Get(path string, description string, fields []field.Field) {
	newRoute := route.Route{
		Protocol:    "GET",
		Path:        a.prefix + path,
		Description: description,
		Fields:      fields,
	}
	a.routes = append(a.routes, newRoute)
}

func (a *Router) Add(rt *Router) {
	a.routers = append(a.routers, rt)
}
