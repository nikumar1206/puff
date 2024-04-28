# Puff

## Vision

- Strictly build upon golang's net/http
- Simplicity where possible

```golang
import "puff/App"

type AppArg struct {
    name string
}
h := AppArg{name: "name_here"}
app := App.New(h)

```

## Features

- Structured Logging
- Routers and nested Routers
- Middlewares
- Auto Open API/Swagger spec generation
- RequestId generation/Tracing Middleware
