package puff

import (
	"log/slog"
	"net/http"
	"reflect"
)

func muxAddHandleFunc(mux *http.ServeMux, route *Route) {
	handler := funcToHandler(route.Handler, route.Parameters, route.Schema)
	mux.Handle(route.Pattern, handler)
}

func handlerToFunc(h http.Handler) HandlerFunc {
	return func(c *Context, _ *struct{}) {
		h.ServeHTTP(c.ResponseWriter, c.Request)
	}
}

func funcToHandler(handler HandlerFunc, parameters map[string]Param, schema reflect.Type) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := NewContext(w, *r)
		params, err := NewParametersInterface(r, schema, parameters)
		if err != nil {
			writeErrorResponse(w, 500, "Creating parameters interface failed with: "+err.Error())
			slog.Info("Creating parameters interface failed with: ", err.Error(), "")
			panic(err.Error())
		}
		handler(context, params)
	})
}

type NetHTTPMiddlewareType func(http.Handler) http.Handler

// converts a net http middleware to puff compatible middleware
// func WrapNetHTTPMiddleware(f NetHTTPMiddlewareType) Middleware {
// 	return func(next HandlerFunc) HandlerFunc {
// 		return func(c *Context, i *interface{}) {
// 			nextHandler := funcToHandler(next)
// 			(f)(nextHandler).ServeHTTP(c.ResponseWriter, c.Request)
// 		}
// 	}
// }
