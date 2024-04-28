package router

import (
	field "puff/field"
	request "puff/request"
	route "puff/route"
)

type Router struct {
	Prefix  string //(optional) prefix, all Routes underneath will have paths that start with the prefix automatically
	Routers []*Router
	Routes  []route.Route
	// middlewares []Middleware
}

func (a *Router) GET(path string, description string, fields []field.Field, handleFunc func(request.Request) interface{}) {
	newRoute := route.Route{
		Protocol:    "GET",
		Path:        path,
		Description: description,
		Fields:      fields,
		Handler:     handleFunc,
	}
	a.Routes = append(a.Routes, newRoute)
}
func (a *Router) POST(path string, description string, fields []field.Field, handleFunc func(request.Request) interface{}) {
	newRoute := route.Route{
		Protocol:    "POST",
		Path:        path,
		Description: description,
		Fields:      fields,
		Handler:     handleFunc,
	}
	a.Routes = append(a.Routes, newRoute)
}
func (a *Router) PUT(path string, description string, fields []field.Field, handleFunc func(request.Request) interface{}) {
	newRoute := route.Route{
		Protocol:    "PUT",
		Path:        path,
		Description: description,
		Fields:      fields,
		Handler:     handleFunc,
	}
	a.Routes = append(a.Routes, newRoute)
}
func (a *Router) PATCH(path string, description string, fields []field.Field, handleFunc func(request.Request) interface{}) {
	newRoute := route.Route{
		Protocol:    "POST",
		Path:        path,
		Description: description,
		Fields:      fields,
		Handler:     handleFunc,
	}
	a.Routes = append(a.Routes, newRoute)
}

func (a *Router) Add(rt *Router) {
	a.Routers = append(a.Routers, rt)
}
