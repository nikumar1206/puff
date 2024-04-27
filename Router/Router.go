package router

import (
	field "puff/field"
	route "puff/route"
)

type Router struct {
	Prefix  string //(optional) prefix, all Routes underneath will have paths that start with the prefix automatically
	Routers []*Router
	Routes  []route.Route
	// middlewares []Middleware
}

func (a *Router) Get(path string, description string, fields []field.Field) {
	newRoute := route.Route{
		Protocol:    "GET",
		Path:        path,
		Description: description,
		Fields:      fields,
	}
	a.routes = append(a.routes, newRoute)
}

func (a *Router) Add(rt *Router) {
	a.Routers = append(a.Routers, rt)
}
