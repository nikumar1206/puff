package route

import (
	field "puff/field"
	"puff/request"
)

type Route struct {
	RouterName string
	Protocol   string
	Path       string
	Pattern    string
	// handle      func(*http.Request) //their handle function
	Description string
	Fields      []field.Field // should have a name, type, description
	Handler     func(request.Request) interface{}
	// should probably have responses (200 OK followed by json, 400 Invalid Paramaters, etc...)
}
