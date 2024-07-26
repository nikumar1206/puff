package puff

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
)

type Reference struct {
	Ref         string `json:"$ref"`
	Summary     string `json:"$summary"`
	Description string `json:"$description"`
}

// OpenAPI struct represents the root of the OpenAPI document.
type OpenAPI struct {
	SpecVersion       string                `json:"openapi"`
	Info              Info                  `json:"info"`
	JSONSchemaDialect string                `json:"jsonSchemaDialect"`
	Servers           []Server              `json:"servers"`
	Paths             Paths                 `json:"paths"`
	Webhooks          map[string]any        `json:"webhooks"`
	Components        Components            `json:"components"`
	Security          []SecurityRequirement `json:"security"`
	Tags              []Tag                 `json:"tags"`
	ExternalDocs      ExternalDocumentation `json:"externalDocs"`
}

// Info struct provides metadata about the API.
type Info struct {
	Title          string  `json:"title"`
	Summary        string  `json:"summary"`
	Description    string  `json:"description"`
	TermsOfService string  `json:"termsOfService"`
	Contact        Contact `json:"contact"`
	License        License `json:"license"`
	Version        string  `json:"version"`
}

// Contact struct contains contact information for the API.
type Contact struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Email string `json:"email"`
}

// License struct contains license information for the API.
type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Server struct represents a server object in OpenAPI.
type Server struct {
	URL         string                    `json:"url"`
	Description string                    `json:"description"`
	Variables   map[string]ServerVariable `json:"variables"`
}

// Components struct holds reusable objects for different aspects of the OAS.
type Components struct {
	Schemas         map[string]any `json:"schemas"`
	Responses       map[string]any `json:"responses"`
	Parameters      map[string]any `json:"parameters"`
	Examples        map[string]any `json:"examples"`
	RequestBodies   map[string]any `json:"requestBodies"`
	Headers         map[string]any `json:"headers"`
	SecuritySchemes map[string]any `json:"securitySchemes"`
	Links           map[string]any `json:"links"`
	Callbacks       map[string]any `json:"callbacks"`
	PathItems       map[string]any `json:"pathItems"`
}

// Tag struct represents a tag used by the OpenAPI document.
type Tag struct {
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	ExternalDocs ExternalDocumentation `json:"externalDocs"`
}

// ExternalDocumentation struct provides external documentation for the API.
type ExternalDocumentation struct {
	Description string `json:"description"`
	URL         string `json:"url"`
}

type Paths map[string]PathItem

// PathItem struct describes operations available on a single path.
type PathItem struct {
	Ref         string                 `json:"$ref"`
	Summary     string                 `json:"summary"`
	Description string                 `json:"description"`
	Get         *Operation             `json:"get,omitempty"`
	Put         *Operation             `json:"put,omitempty"`
	Post        *Operation             `json:"post,omitempty"`
	Delete      *Operation             `json:"delete,omitempty"`
	Options     *Operation             `json:"options,omitempty"`
	Head        *Operation             `json:"head,omitempty"`
	Patch       *Operation             `json:"patch,omitempty"`
	Trace       *Operation             `json:"trace,omitempty"`
	Servers     []Server               `json:"servers"`
	Parameters  []ParameterOrReference `json:"parameters"`
}

type SecurityRequirement map[string][]string

// Operation struct describes an operation in a PathItem.
type Operation struct {
	Tags         []string               `json:"tags"`
	Summary      string                 `json:"summary"`
	Description  string                 `json:"description"`
	ExternalDocs ExternalDocumentation  `json:"externalDocs"`
	OperationID  string                 `json:"operationId"`
	Parameters   []ParameterOrReference `json:"parameters"`
	RequestBody  RequestBodyOrReference `json:"requestBody"`
	Responses    map[string]Response    `json:"responses"`
	Callbacks    map[string]Callback    `json:"callbacks"`
	Deprecated   bool                   `json:"deprecated"`
	Security     []SecurityRequirement  `json:"security"`
	Servers      []Server               `json:"servers"`
}

// ParameterOrReference is a union type representing either a Parameter Object or a Reference Object.
type ParameterOrReference struct {
	Reference string    `json:"$ref,omitempty"`
	Parameter Parameter `json:"-"`
}

// Parameter struct describes a parameter in OpenAPI.
type Parameter struct {
	Name            string `json:"name"`
	In              string `json:"in"`
	Description     string `json:"description"`
	Required        bool   `json:"required"`
	Deprecated      bool   `json:"deprecated"`
	AllowEmptyValue bool   `json:"allowEmptyValue"`
	Style           string `json:"style"`
	Explode         bool   `json:"explode"`
	AllowReserved   bool   `json:"allowReserved"`
	Schema          Schema `json:"schema"`
}

// RequestBodyOrReference is a union type representing either a Request Body Object or a Reference Object.
type RequestBodyOrReference struct {
	Reference   string      `json:"$ref,omitempty"`
	RequestBody RequestBody `json:"-"`
}

// RequestBody struct describes a request body in OpenAPI.
type RequestBody struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content"`
	Required    bool                 `json:"required"`
}

// MediaType struct describes a media type object in OpenAPI.
type MediaType struct {
	Schema   Schema         `json:"schema"`
	Example  any            `json:"example"`
	Examples map[string]any `json:"examples"`
}

// Schema struct represents a schema object in OpenAPI.
type Schema struct {
	// Define your schema fields based on your specific requirements
	// Example fields could include type, format, properties, etc.
	// This can be expanded based on the needs of your application.
}

// OpenAPIResponse struct describes possible responses in OpenAPI.
type OpenAPIResponse struct {
	Description string               `json:"description"`
	Headers     map[string]Header    `json:"headers,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty"`
	Links       map[string]Link      `json:"links,omitempty"`
}

type Callback map[string]PathItem

type Example struct {
	Summary       string `json:"summary,omitempty"`
	Description   string `json:"description,omitempty"`
	Value         any    `json:"value,omitempty"`
	ExternalValue string `json:"externalValue,omitempty"`
}

type Header struct {
	// Add header fields as needed
}

type Link struct {
	OperationRef string `json:"operationRef,omitempty"`
	OperationID  string `json:"operationId,omitempty"`
	Parameters   any    `json:"parameters,omitempty"`
	RequestBody  any    `json:"requestBody,omitempty"`
	Description  string `json:"description,omitempty"`
	Server       Server `json:"server,omitempty"`
}

type Encoding struct {
	ContentType   string            `json:"contentType,omitempty"`
	Headers       map[string]Header `json:"headers,omitempty"`
	Style         string            `json:"style,omitempty"`
	Explode       bool              `json:"explode,omitempty"`
	AllowReserved bool              `json:"allowReserved,omitempty"`
}

type ServerVariable struct {
	Enum        []string `json:"enum,omitempty"`
	Default     string   `json:"default"`
	Description string   `json:"description,omitempty"`
}

func GenerateOpenAPIUI(document, title, docsURL string) string {
	return fmt.Sprintf(openAPIHTML, title, docsURL)
}

func addRoute(router Router, route Route, tags *[]Tag, tagNames *[]string, paths *Paths) {
	tag := router.Tag

	if tag == "" {
		tag = router.Name
	}
	if !slices.Contains(*tagNames, tag) {

		*tagNames = append(*tagNames, tag)
		*tags = append(*tags, Tag{Name: tag, Description: ""})
	}

	description := "This route does"
	summary := description
	if len(summary) > 100 {
		summary = summary[:100] + " ..."
	}
	pathMethod := &Operation{
		Summary:     summary,
		OperationID: "",
		Tags:        []string{tag},
		Parameters:  []ParameterOrReference{},
		Responses:   map[string]Response{},
		Description: description, // TODO: needs to be dynamic on route
	}
	pathItem := (*paths)[route.fullPath]
	switch route.Protocol {
	// TODO: handle other protocols
	case http.MethodGet:
		pathItem.Get = pathMethod
	case http.MethodPut:
		pathItem.Put = pathMethod
	case http.MethodPatch:
		pathItem.Patch = pathMethod
	case http.MethodPost:
		pathItem.Post = pathMethod
	case http.MethodDelete:
		pathItem.Delete = pathMethod
	}
	(*paths)[route.fullPath] = pathItem

}

func GenerateOpenAPISpec(
	appName string,
	appVersion string,
	rootRouter Router,
) (string, error) {
	var tags []Tag
	var tagNames []string
	var paths = make(Paths)
	for _, route := range rootRouter.Routes {
		addRoute(rootRouter, *route, &tags, &tagNames, &paths)
	}
	for _, router := range rootRouter.Routers {
		for _, route := range router.Routes {
			addRoute(*router, *route, &tags, &tagNames, &paths)
		}
	}

	info := Info{
		Version: appVersion,
		Title:   appName,
	}
	openapi := OpenAPI{
		SpecVersion: "3.1.0",
		Info:        info,
		Servers:     []Server{},
		Tags:        tags,
		// FIX ME: SERVERS SHOULD BE SPECIFIED IN THE APP CONFIGURATION
		// FIX ME: THE DEFAULT SERVER SHOULD BE THE NETWORK IP: PORT
		Paths: paths,
	}
	openapiJSON, err := json.Marshal(openapi)
	if err != nil {
		return "", err
	}
	return string(openapiJSON), nil
}
