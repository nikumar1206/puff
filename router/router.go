package router

import (
	field "puff/field"
	request "puff/request"
	route "puff/route"
)

type Router struct {
	Name    string
	Prefix  string //(optional) prefix, all Routes underneath will have paths that start with the prefix automatically
	Routers []*Router
	Routes  []route.Route
	// middlewares []Middleware
}

func (a *Router) GET(path string, description string, fields []field.Field, handleFunc func(request.Request) interface{}) {
	newRoute := route.Route{
		RouterName:  a.Name,
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
		RouterName:  a.Name,
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
		RouterName:  a.Name,
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
		RouterName:  a.Name,
		Protocol:    "POST",
		Path:        path,
		Description: description,
		Fields:      fields,
		Handler:     handleFunc,
	}
	a.Routes = append(a.Routes, newRoute)
}

func (a *Router) IncludeRouter(rt *Router) {
	a.Routers = append(a.Routers, rt)
}
