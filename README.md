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

- Separate Makefile build commands. Currently everything running via `make reload`
- Fix route collision issue
- Server Side Events do not support all event types
- Improve documentation everywhere
- Better adherance to OpenAPI spec
- Change the name for IncludeRouter/IncludeMiddleware. It is obtuse.
- JSON should only do JSON logging. Text logger, shouldn't indent JSON
- Add a Middleware skipper function that depends on context.
- Remove color package dependency.
- Allow configuring the logger and making it more generic
  - Allow indenting/non-indenting in JSON logger.
- Should puff.Context have context on the app? So that it knows to log in JSON or not.
- find something better than regex matching
