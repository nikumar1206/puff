package puff

import "net/http"

type Context interface {
	// used to serialize the provided response appropriately.
	SendResponse() func(Response) error
	// returns the Request interface
	GetRequest() func() *http.Request
	// returns the response writer for the route
	GetRequestWriter() func() *http.ResponseWriter
	// sets the status code of the current route
	SetStatusCode() func(int)
	// returns the status code of the current route
	GetStatusCode() func() int
}



func (c *Context) GetRequest() {

}
