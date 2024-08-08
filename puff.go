// Package puff provides primitives for implementing a Puff Server
package puff

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
	// DocsReload, if true, enables automatic reload on the Swagger documentation page.
	DocsReload bool
	// TLSPublicCertFile specifies the file for the TLS certificate (usually .pem or .crt).
	TLSPublicCertFile string
	// TLSPrivateKeyFile specifies the file for the TLS private key (usually .key).
	TLSPrivateKeyFile string
}

func App(c *AppConfig) *PuffApp {
	r := &Router{Name: "Puff Default", Tag: "Default", Description: "Puff Default Router"}
	if c.Version == "" {
		c.Version = "0.0.0"
	}

	return &PuffApp{
		Name:              c.Name,
		Version:           c.Version,
		DocsURL:           c.DocsURL,
		DocsReload:        c.DocsReload,
		TLSPublicCertFile: c.TLSPublicCertFile,
		TLSPrivateKeyFile: c.TLSPrivateKeyFile,
		RootRouter:        r,
	}
}

func DefaultApp(name string) *PuffApp {
	c := AppConfig{
		Version: "1.0.0",
		Name:    name,
		DocsURL: "/docs",
	}
	a := App(&c)
	a.Logger = DefaultLogger()
	return a
}
