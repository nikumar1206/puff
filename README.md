# Puff

An extensible, performant, and reliable framework inspired by FastAPI.

![Alt](https://repobeats.axiom.co/api/embed/66ccd66540fab2ca27806fc48acba71ab93721d5.svg "Repobeats analytics image")

## Features

- Automatic OpenAPI documentation generation and hosting.
- Tree-structure style Routers to group APIs.
- Extensible middlewares.
- Customizable logger with structured and prettier logging.
- Adhere to standards set by net/http, and RFC-compliant.
- Simplicity where possible and build upon the goated stdlib when possible.
  - Only has 2 external dependencies!

## Quickstart

### Installation

```bash
go get -u github.com/nikumar1206/puff
```

Creating a new server using Puff is simple:

```golang
import "puff"

app := puff.DefaultApp("Cool Demo App")
app.Use(middleware.Tracing())
app.Use(middleware.Logging())

app.Get("/health", nil, func(c *puff.Context) {
		c.SendResponse(puff.GenericResponse{Content: "Server is Healthy!"})
})

app.ListenAndServe(":8000")
```

This will begin the application and serve the docs at `http://localhost:8000/docs`.

The DefaultApp sets up some great defaults, but you can specify a custom config with `puff.App()`.

We also recommend setting up your own logger:

```golang
app.Logger = puff.NewLogger(puff.LoggerConfig{
		Level:      slog.LevelDebug,
		Colorize:   true,
		UseJSON:    false,
		TimeFormat: time.DateTime,
})
```

### RoadMap

##### [View the roadmap](./roadmap.md)

##### TODO: improve docs with example on creating a simple router, and contribution docs.
