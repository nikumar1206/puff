package route

import (
	field "puff/field"
)

type Route struct {
	Protocol string
	Path     string
	// handle      func(*http.Request) //their handle function
	Description string
	Fields      []field.Field // should have a name, type, description
	// should probably have responses (200 OK followed by json, 400 Invalid Paramaters, etc...)
}
