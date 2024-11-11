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

	if app.Name != "TestApp" {
		t.Errorf("Expected app name 'TestApp', got '%s'", app.Name)
	}
	if app.Version != "1.2.3" {
		t.Errorf("Expected version '1.2.3', got '%s'", app.Version)
	}
	if app.DocsURL != "/test-docs" {
		t.Errorf("Expected DocsURL '/test-docs', got '%s'", app.DocsURL)
	}
	if app.TLSPublicCertFile != "cert.pem" {
		t.Errorf("Expected TLSPublicCertFile 'cert.pem', got '%s'", app.TLSPublicCertFile)
	}
	if app.TLSPrivateKeyFile != "key.pem" {
		t.Errorf("Expected TLSPrivateKeyFile 'key.pem', got '%s'", app.TLSPrivateKeyFile)
	}
	if app.OpenAPI != nil {
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

	if app.Version != "0.0.0" {
		t.Errorf("Expected default version '0.0.0', got '%s'", app.Version)
	}
}

func TestDefaultApp(t *testing.T) {
	app := DefaultApp("DefaultAppTest")

	if app.Name != "DefaultAppTest" {
		t.Errorf("Expected app name 'DefaultAppTest', got '%s'", app.Name)
	}
	if app.Version != "1.0.0" {
		t.Errorf("Expected default version '1.0.0', got '%s'", app.Version)
	}
	if app.DocsURL != "/docs" {
		t.Errorf("Expected default DocsURL '/docs', got '%s'", app.DocsURL)
	}
	if app.Logger == nil {
		t.Fatalf("Expected Logger to be initialized")
	}
}
