// Package puff provides primitives for implementing a Puff Server
package puff

import "log/slog"

type HandlerFunc func(*Context)
type Middleware func(next HandlerFunc) HandlerFunc

// AppConfig is defines PuffApp parameters.
type AppConfig struct {
	// Name is the application name
	Name string
	// Version is the application version.
	Version string
	// DocsURL is the Router prefix for Swagger documentation. Can be "" to disable Swagger documentation.
	DocsURL string
	// TLSPublicCertFile specifies the file for the TLS certificate (usually .pem or .crt).
	TLSPublicCertFile string
	// TLSPrivateKeyFile specifies the file for the TLS private key (usually .key).
	TLSPrivateKeyFile string
	// OpenAPI configuration. Gives users access to the OpenAPI spec generated. Can be manipulated by the user.
	OpenAPI *OpenAPI
	// SwaggerUIConfig is the UI specific configuration.
	SwaggerUIConfig *SwaggerUIConfig
}

func App(c *AppConfig) *PuffApp {
	r := &Router{Name: "Default", Tag: "Default", Description: "Default Router"}
	if c.Version == "" {
		c.Version = "0.0.0"
	}

	a := &PuffApp{
		Name:              c.Name,
		Version:           c.Version,
		DocsURL:           c.DocsURL,
		TLSPublicCertFile: c.TLSPublicCertFile,
		TLSPrivateKeyFile: c.TLSPrivateKeyFile,
		RootRouter:        r,
		OpenAPI:           c.OpenAPI,
	}
	a.RootRouter.puff = a
	a.RootRouter.Responses = Responses{}
	return a
}

func DefaultApp(name string) *PuffApp {
	c := AppConfig{
		Version: "1.0.0",
		Name:    name,
		DocsURL: "/docs",
	}
	a := App(&c)
	a.Logger = DefaultLogger()
	slog.SetDefault(a.Logger)
	return a
}
