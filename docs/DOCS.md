# Creating a Puff App

```golang
package main
import "puff"

func main(){
    app *puff.PuffApp := puff.DefaultApp()
}
```

More options are available when using puff.App and puff.Config.

```golang
config := puff.Config{
    Name: "my_app_name_here",
    Version: "v0.0.1",
    DocsURL: "/docs", // DocsURL is the page to serve the OpenAPI Schema.
}
app *puff.PuffApp := puff.App(config)
```

# Creating a Router

```golang
package main
import "puff"

func main(){
    app *puff.PuffApp := puff.DefaultApp()
    router *puff.Router := puff.NewRouter(
        Name: "my_router_name_here",
        Prefix: ""
    )
    app.IncludeRouter(router) // Including the router is mandatory to serve the router on the app.
}
```

It is possible to do `router.IncludeRouter(anotherRouter)`.

# Example Router Tree

<img src="example router structure.png"></img>

# Writing a GET Request

```golang
router.GET("/", puff.Field{}, func(c *puff.Context){
    c.SendResponse(puff.GenericResponse{Content: "Hello, world!"})
})
```

# Response Types

There are a few provided response types that you can send through `*puff.Context.SendResponse` during route handling.

## JSONResponse

```golang
package main
import "puff"

type User struct {
    Name string `json:"name"`
}

func main(){
    app := puff.DefaultApp()
    app.Get("/", "", nil, func (c *puff.Context){
        user1 := User{
            Name: "John Doe",
        }
        c.SendResponse(puff.JSONResponse{
            StatusCode: 200, // StatusCode defaults to 200 if not provided.
            Content: user1,
        })
    })
```

## HTMLResponse

```golang
c.SendResponse(puff.HTMLResponse{
    StatusCode: 200, // StatusCode defaults to 200 if not provided.
    Content: `<pre>Hello, world!</pre>`, // Content defaults to an empty string if not provided.
})
```

## FileResponse

```golang
c.SendResponse(puff.FileResponse{
        StatusCode:  200 // StatusCode defaults to 200 if not provided.
	FilePath:    "path/to/assets/image.jpg" // FilePath to file
	FileContent []byte
	ContentType string // ContentType is inferred from extenstion of
})
```

## StreamingResponse

```golang
c.SendResponse(puff.StreamingResponse{
    StatusCode: 200,
    StreamHandler: func(s *chan string){
        defer close(s) // The connection does not close until you close the channel.
        for i := range 3 {
            s <- i
            time.Sleep(5 * time.Second)
        }
    }
})
```

## RedirectResponse

```golang
c.SendResponse(puff.RedirectResponse{
    StatusCode: 308,
    To: "https://google.com"
})
```

## GenericResponse

```golang
c.SendResponse(puff.GenericResponse{
    StatusCode: 200,
    Content: "Hello, world!"
    ContentType: "text/plain"
})
```

# Input Schemas

Input schemas specify what types of inputs your application takes.

Example Usage:

```
package main
import "puff"

type HelloWorldInput struct {
    Name string `kind:"query"`
}

func main(){
    app := puff.App("Input Schemas Example")

    hello_world_input := new(HelloWorldInput)
    app.Get(path: "/", description: "greets you by name", fields: hello_world_input, func (c *Context) {
        c.SendResponse(hello_world_input.Name)
    })
}
```

The schema, `HelloWorldInput` in this example, specifies a query parameter of type string.

The struct tag can take:
| Field | Required | Description | Possible Values |
| -------- | -- | -- |------- |
| kind | Yes | where should the parameter be found | `query`, `path`, `header`, `cookie` |
| description | No | a brief description of the parameter | anything |
| required | No | specifies if the field is required. defaults to true for everything except cookie | `true`, `false`|
| deprecated | No | marks field as deprecated. defaults to false. | `true`, `false`|
| format | No | the format of the parameter. | examples: `email`, `password`, `uint64`|

When passing in the input, it must be a pointer to something with the input schema as the type.

**NOTE**: The proccessing of the input schema may panic. Here's what the panic messages mean.

| Message                                                | Meaning                                                        |
| ------------------------------------------------------ | -------------------------------------------------------------- |
| field must be POINTER to structure                     | The value you passed in for the input schema is not a pointer. |
| field must be pointer to STRUCT                        | The value you passed in for the input schma is not a struct.   |
| type of field \_\_ must be string or int               | The type of the field must be string or int.                   |
| specified kind on field \_\_ in struct tag must be ... | The kind on the field is not a supported kind.                 |

# Middlewares

Middlewares provide many useful tools to enhance your application. Puff comes with many middlewares in the middleware package.

To install the middleware package:

`go get https://github.com/nikumar1206/puff/middleware`

## Attaching a Middleware

```golang
package main

import (
	"github.com/nikumar1206/puff"
	"github.com/nikumar1206/puff/middleware"
)

func main() {
	app := puff.DefaultApp("Restaurant Microservice")
	app.Use(middleware.CSRF())
}
```

The middleware package provides many middlewares. You can view the middleware docs at [the middleware pkg documentation](https://pkg.go.dev/github.com/nikumar1206/puff/middleware#section-documentation).

## The Middleware Standard

Each middleware should have all the following.

### Configuration

```golang
type MyMiddlewareConfig struct {
    // include configuration fields here
}
```

The `MyMiddlewareConfig` should contain configuration fields for your middleware.

### Default Configuration

```golang
var DefaultMyMiddlewareConfig MyMiddlewareConfig = MyMiddlewareConfig{
    // specify default configuration here
}
```

The `DefaultMyMiddlewareConfig` should completely populate MyMiddlewareConfig.

### Middleware Function

```golang
func MyMiddleware() puff.Middleware
```

The `MyMiddleware` function should return a `puff.Middleware` without asking for a configuration. Instead, it should use the default configuration at `DefaultMyMiddlewareConfig`.

### Middleware Function with Configuration

```golang
func MyMiddlewareWithConfig(*MyMiddlewareConfig) puff.Middleware
```

The `MyMiddlewareWithConfig` function should take in a pointer to `MyMiddlewareConfig` and return a `puff.Middleware`.
