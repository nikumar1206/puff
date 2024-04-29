package route

import "github.com/nikumar1206/puff/request"

type Route struct {
	Protocol string
	Pattern  string
	Path     string
	Handler  func(request.Request) interface{}
	// should probably have responses (200 OK followed by json, 400 Invalid Paramaters, etc...)
}
