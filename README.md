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

## RoadMap

- Logging
  - Allow users to use their own time format.
- Middlewares
  - Some sort of Session Middleware/Authentication Middleware? Potentially both?
  - Panic Handler? (Wrap panics into 500) ✅
  - Allow configuration of middleware settings. When adding middlewares, they can pass a config. Current middleware style will be default.
  - Allow attaching middlewares to routers (far future)
- Fixes
  - Remove router name and add "tags" instead
    - Also believe this is broken.
    - Router name for Drinks Router doesnt appear.
  - Fix the way 'description' is set.



## Definitely need to fix/improve
- Running from just `go run examples/restaurant_app/main.go` doesnt work. We need to fix
