package puff

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

// Context provides functionality for the route.
type Context struct {
	// Request is the underlying *http.Request object.
	Request *http.Request
	// ResponseWriter is the underlying http.ResponseWriter object.
	ResponseWriter http.ResponseWriter
	// WebSocket represents WebSocket connection and its related context, connection, and events.
	// WebSocket will be nil if the route does not use websockets.
	WebSocket  *WebSocket
	statusCode int
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Request:        r,
		ResponseWriter: w,
	}
}

func (ctx *Context) isWebSocket() bool {
	return ctx.GetHeader("Upgrade") == "websocket" &&
		ctx.GetHeader("Connection") == "Upgrade" &&
		ctx.GetHeader("Sec-WebSocket-Version") == "13"
}

// GetHeader gets the value of a request header with key k.
// It returns an empty string if not found.
func (ctx *Context) GetHeader(k string) string {
	return ctx.Request.Header.Get(k)
}

// GetHeader gets the value of a response header with key k.
// It returns an empty string if not found.
func (ctx *Context) GetResponseHeader(k string) string {
	return ctx.ResponseWriter.Header().Get(k)
}

// SetHeader sets the value of the response header k to v.
func (ctx *Context) SetHeader(k, v string) {
	ctx.ResponseWriter.Header().Set(k, v)
}

// SetContentType sets the content type of the response.
func (ctx *Context) SetContentType(v string) {
	ctx.SetHeader("Content-Type", v)
}

// SetStatusCode sets the status code of the response.
func (ctx *Context) SetStatusCode(sc int) {
	ctx.ResponseWriter.WriteHeader(sc)
	ctx.statusCode = sc
}

// GetStatusCode returns the status code. If response not written, returns default 0.
func (ctx *Context) GetStatusCode() int {
	return ctx.statusCode
}

// below are methods that are more utility focused.

// GetRequestID gets the X-Request-ID if set (empty string if not set).
// puff/middleware provides a tracing middleware the sets X-Request-ID.
func (ctx *Context) GetRequestID() string {
	return ctx.GetResponseHeader("X-Request-ID")
}

// SendResponse sends res back to the client.
// Any errors at this point will be logged and the request will fail.
func (c *Context) SendResponse(res Response) {
	c.SetContentType(res.GetContentType())
	c.SetStatusCode(res.GetStatusCode())
	err := res.WriteContent(c)
	if err != nil {
		msg := fmt.Sprintf(
			"[%s] An unexpected error occured while writing content with context: %s.",
			c.GetRequestID(),
			err.Error(),
		)
		slog.Error(msg)
		fmt.Fprint(c.ResponseWriter, "An unknown error occured.")
	}
}

// GetBearerToken gets the Bearer token if it exists.
// This will work if the request contains an Authorization header
// that has this syntax: Bearer this_token_here.
func (ctx *Context) GetBearerToken() string {
	bt := ctx.GetHeader("Authorization")

	token_arr := strings.Split(bt, "Bearer ")

	if len(token_arr) > 1 {
		return token_arr[1]
	}

	return ""
}

// below are methods that are more error message focused.

func (ctx *Context) response(status_code int, message string, a ...any) {
	ctx.SendResponse(JSONResponse{
		StatusCode: status_code,
		Content: map[string]any{
			"error": fmt.Sprintf(message, a...),
		},
	})
}

// BadRequest returns a json response with status code 400
// a key error and a value of the formatted string from
// message and the arguments following.
func (ctx *Context) BadRequest(message string, a ...any) {
	ctx.response(400, message, a...)
}

// BadRequest returns a json response with status code 500
// with a key error and a value of the formatted string from
// message and the arguments following.
func (ctx *Context) InternalServerError(message string, a ...any) {
	ctx.response(500, message, a...)
}
