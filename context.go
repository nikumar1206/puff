package puff

import (
	"log/slog"
	"net/http"
	"strings"
)

type Context struct {
	// original http.request object
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	WebSocket      WebSocket //WebSocket (if valid websocket connection)
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	slog.Debug("Initiating new context for request")
	return &Context{
		Request:        r,
		ResponseWriter: w,
	}
}

// returns "" if provided key cannot be found
func (ctx *Context) GetHeader(k string) string {
	return ctx.Request.Header.Get(k)
}

// sets a response header
func (ctx *Context) SetHeader(k, v string) {
	ctx.ResponseWriter.Header().Set(k, v)
}

func (ctx *Context) SetContentType(v string) {
	ctx.SetHeader("Content-Type", v)
}

// sets the respons status code
func (ctx *Context) SetStatusCode(sc int) {
	ctx.ResponseWriter.WriteHeader(sc)
}

// below are methods that are more utility focused.

// provides x-request-id from headers if set, else returns ""
func (ctx *Context) GetRequestID() string {
	return ctx.GetHeader("X-Request-ID")
}

func (ctx *Context) SendResponse(res Response) {
	res.WriteContent(ctx.ResponseWriter, ctx.Request)
}

// will try to return bearer token if exists, else returns ""
func (ctx *Context) GetBearerToken() string {
	bt := ctx.GetHeader("Authorization")

	token_arr := strings.Split(bt, "Bearer ")

	if len(token_arr) > 1 {
		return token_arr[1]
	}

	return ""
}
