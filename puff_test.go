package puff

import (
	"testing"
)

func TestApp(t *testing.T) {
	// Test with all configuration fields set
	config := AppConfig{
		Name:              "TestApp",
		Version:           "1.2.3",
		DocsURL:           "/test-docs",
		TLSPublicCertFile: "cert.pem",
		TLSPrivateKeyFile: "key.pem",
	}
	app := App(&config)

	if app.Config.Name != "TestApp" {
		t.Errorf("Expected app name 'TestApp', got '%s'", app.Config.Name)
	}
	if app.Config.Version != "1.2.3" {
		t.Errorf("Expected version '1.2.3', got '%s'", app.Config.Version)
	}
	if app.Config.DocsURL != "/test-docs" {
		t.Errorf("Expected DocsURL '/test-docs', got '%s'", app.Config.DocsURL)
	}
	if app.Config.TLSPublicCertFile != "cert.pem" {
		t.Errorf("Expected TLSPublicCertFile 'cert.pem', got '%s'", app.Config.TLSPublicCertFile)
	}
	if app.Config.TLSPrivateKeyFile != "key.pem" {
		t.Errorf("Expected TLSPrivateKeyFile 'key.pem', got '%s'", app.Config.TLSPrivateKeyFile)
	}
	if app.Config.OpenAPI != nil {
		t.Errorf("Expected OpenAPI to not be set.")
	}
	if app.RootRouter == nil {
		t.Fatalf("Expected RootRouter to be initialized")
	}
	if app.RootRouter.Name != "Default" {
		t.Errorf("Expected RootRouter name 'Default', got '%s'", app.RootRouter.Name)
	}
	if app.RootRouter.Tag != "Default" {
		t.Errorf("Expected RootRouter tag 'Default', got '%s'", app.RootRouter.Tag)
	}
}

func TestApp_DefaultVersion(t *testing.T) {
	// Test with default version when version is not provided in the config
	config := AppConfig{
		Name: "TestAppWithDefaultVersion",
	}
	app := App(&config)

	if app.Config.Version != "0.0.0" {
		t.Errorf("Expected default version '0.0.0', got '%s'", app.Config.Version)
	}
}

func TestDefaultApp(t *testing.T) {
	app := DefaultApp("DefaultAppTest")

	if app.Config.Name != "DefaultAppTest" {
		t.Errorf("Expected app name 'DefaultAppTest', got '%s'", app.Config.Name)
	}
	if app.Config.Version != "1.0.0" {
		t.Errorf("Expected default version '1.0.0', got '%s'", app.Config.Version)
	}
	if app.Config.DocsURL != "/docs" {
		t.Errorf("Expected default DocsURL '/docs', got '%s'", app.Config.DocsURL)
	}
}
