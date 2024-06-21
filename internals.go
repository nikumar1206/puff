package puff

import "net/http"

func muxAddHandleFunc(mux *http.ServeMux, route *Route) {
	handler := funcToHandler(route.Handler)
	mux.Handle(route.Pattern, handler)
}

func handlerToFunc(h http.Handler) HandlerFunc {
	return func(c *Context) {
		h.ServeHTTP(c.ResponseWriter, c.Request)
	}
}

func funcToHandler(f HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := NewContext(w, *r)
		f(c)
	})
}

type NetHTTPMiddlewareType func(http.Handler) http.Handler

// converts a net http middleware to puff compatible middleware
func WrapNetHTTPMiddleware(f NetHTTPMiddlewareType) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			nextHandler := funcToHandler(next)
			(f)(nextHandler).ServeHTTP(c.ResponseWriter, c.Request)
		}
	}

}
