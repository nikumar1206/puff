package puff

// import (
// 	"net/http"
// )

// func muxAddHandleFunc(server http.Server, route *Route) {
// 	handler := puffHandlerFuncToHTTPHandler(route.Handler)
// 	mux.Handle(route.Pattern, handler)
// 	server.Handler.ServeHTTP(http.ResponseWriter, *http.Request)
// }

// // func httpHandlerToPuffHandlerFunc(h http.Handler) HandlerFunc {
// // 	return func(c *Context) {
// // 		h.ServeHTTP(c.ResponseWriter, c.Request)
// // 	}
// // }

// func puffHandlerFuncToHTTPHandler(f HandlerFunc) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		c := NewContext(w, *r)
// 		f(c)
// 	})
// }

// type NetHTTPMiddlewareType func(http.Handler) http.Handler

// // converts a net http middleware to puff compatible middleware
// func WrapNetHTTPMiddleware(f NetHTTPMiddlewareType) Middleware {
// 	return func(next HandlerFunc) HandlerFunc {
// 		return func(c *Context) {
// 			nextHandler := puffHandlerFuncToHTTPHandler(next)
// 			(f)(nextHandler).ServeHTTP(c.ResponseWriter, c.Request)
// 		}
// 	}

// }
