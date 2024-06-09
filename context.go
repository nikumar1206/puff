package puff

import (
	"net/http"
	"strings"
)

type Context struct {
	// original http.request object
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

func NewContext(w http.ResponseWriter, r http.Request) *Context {
	return &Context{
		Request:        &r,
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
	switch r := res.(type) {
	case JSONResponse:
		handleJSONResponse(ctx.ResponseWriter, ctx.Request, r)
	case HTMLResponse:
		handleHTMLResponse(ctx.ResponseWriter, ctx.Request, r)
	case FileResponse:
		handleFileResponse(ctx.ResponseWriter, ctx.Request, r)
	case StreamingResponse:
		handleStreamingResponse(ctx.ResponseWriter, r)
	case GenericResponse:
		handleGenericResponse(ctx.ResponseWriter, ctx.Request, r)
	default:
		writeErrorResponse(ctx.ResponseWriter, http.StatusInternalServerError, "Invalid response type")
	}
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
