package puff

import "net/http"

type Context interface {
	// used to serialize the provided response appropriately.
	SendResponse() func(Response) error
	// returns the Request interface
	Request() func() *http.Request
	// returns the response writer for the route
	RequestWriter() func() *http.ResponseWriter
	// sets the status code of the current route
	SetStatusCode() func(int)
	// returns the status code of the current route
	GetStatusCode() func() int
	GetHeader() func(string) string
	SetHeader(string, string)
	GetHost() string
}

type context struct {
	request *http.Request
	rw      *http.ResponseWriter
	sc      int
}

func (ctx context) GetHost() string {
	return ctx.request.Host
}

// returns "" if provided key cannot be found
func (ctx context) GetHeader(k string) string {
	return ctx.request.Header.Get(k)
}

// sets a response header
func (ctx context) SetHeader(k, v string) {
	(*ctx.rw).Header().Add(k, v)
}

// sets a response header
func (ctx context) GetStatusCode() int {
	return ctx.sc
}

func (ctx context) SetStatusCode(sc int) {
	(*ctx.rw).WriteHeader(sc)
	ctx.sc = sc
}

// provides x-request-id from headers if set
func (ctx context) GetRequestID() string {
	return ctx.GetHeader("X-Request-ID")
}

// func (c *context) writeContentType(value string) {
// 	header := c.Response().Header()
// 	if header.Get(HeaderContentType) == "" {
// 		header.Set(HeaderContentType, value)
// 	}
// }
// func (c *Context) GetRequest() {

// }
