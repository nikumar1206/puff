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
    app.Get("/", puff.Field{}, func (c *puff.Context){
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
