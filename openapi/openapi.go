package openapi

import (
	"encoding/json"
	"fmt"
	"puff/route"
	"slices"
	"strings"
)

type License struct {
	Name string `json:"name"` //MIT, CC-BY-0, etc.
	Url  string `json:"url"`
}
type Info struct {
	Version string  `json:"version"` //ex. 1.0.0
	Title   string  `json:"title"`
	License License `json:"license"`
	// add licensing here
}
type Server struct {
	Url         string `json:"url"`
	Description string `json:"description"`
}
type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

//	type Schema struct {
//		Type string `json:"type"`
//	}
type Parameter struct {
	Name        string `json:"url"`
	In          string `json:"in"` //path, query, header, cookie
	Description string `json:"description"`
	Required    bool   `json:"required"`
	// Schema      Schema `json:"schema"`
	Deprecated bool `json:"deprecated"`
}
type Response struct {
	Description string
	// Headers []Header
	// Content []Content
	//Fix Me: https://swagger.io/specification/#responses-object
}
type Get struct {
	*Method `json:"get"`
}
type Post struct {
	*Method `json:"post"`
}
type Put struct {
	*Method `json:"put"`
}
type Patch struct {
	*Method `json:"patch"`
}
type Method struct {
	// https://swagger.io/specification/#:~:text=style%3A%20simple-,Operation%20Object,-Describes%20a%20single
	Summary     string              `json:"summary"`
	OperationID string              `json:"operationId"`
	Tags        []string            `json:"tags"`
	Parameters  []Parameter         `json:"parameters"`
	Description string              `json:"description"`
	Responses   map[string]Response `json:"responses"`
	Deprecated  bool                `json:"deprecated"`
	//FIX ME: Request-Body
}

type OpenAPI struct {
	SpecVersion string                            `json:"openapi"` //this is the version, should be 3.1.0
	Info        Info                              `json:"info"`
	Servers     []Server                          `json:"servers"`
	Paths       map[string]map[string]interface{} `json:"paths"` //the string key is the path
}

var OPENAPI_UI string = `
<!DOCTYPE html>
<html >
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="theme-color" content="#000000" />
    <meta name="description" content="SwaggerUIMultifold" />
    <link rel="stylesheet" href="//unpkg.com/swagger-editor@5.0.0-alpha.86/dist/swagger-editor.css" />
    <title>%s</title>
  </head>
  <body style="margin:0; padding:0;">
  <div style="position: absolute; top: 0; left:0; width: 100vw; height: 20px; background:black;"></div>
    <section id="swagger-ui"></section>
    <script src="//unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js"></script>
    <script src="//unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-standalone-preset.js"></script>
    <script>
      ui = SwaggerUIBundle({});
      // expose SwaggerUI React globally for SwaggerEditor to use
      window.React = ui.React;
    </script>
    <script src="//unpkg.com/swagger-editor@5.0.0-alpha.86/dist/umd/swagger-editor.js"></script>
    <script>
      SwaggerUIBundle({
        url: '/api/docs/docs.json',
        dom_id: '#swagger-ui',
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset,
        ],
        plugins: [
          SwaggerEditor.plugins.EditorContentType,
          SwaggerEditor.plugins.EditorPreviewAsyncAPI,
          SwaggerEditor.plugins.EditorPreviewApiDesignSystems,
          SwaggerEditor.plugins.SwaggerUIAdapter,
          SwaggerUIBundle.plugins.DownloadUrl,
        ],
        layout: 'StandaloneLayout',
      });
    document.body.onload = function (){
    const sc = document.getElementsByClassName("swagger-container")
    const tb = document.getElementsByClassName("topbar")
    sc[0].removeChild(tb[0])
  }
  </script>
  </body>
</html>
`

func GenerateOpenAPIUI(document string, title string) string {
	return fmt.Sprintf(OPENAPI_UI, title)
}

func GenerateOpenAPISpec(appName string, appVersion string, appRoutes []*route.Route) (string, error) {
	var tags []Tag
	var tagNames []string
	var paths map[string]map[string]interface{} = make(map[string]map[string]interface{})
	for _, r := range appRoutes {
		if !slices.Contains(tagNames, r.RouterName) {
			tagNames = append(tagNames, r.RouterName)
			tags = append(tags, Tag{Name: r.RouterName, Description: "replace this: line 138 openapi.go"})
		}
		pathMethod := Method{
			Summary:     "",
			OperationID: "",
			Tags:        []string{r.RouterName},
			Parameters:  []Parameter{},
			Responses:   map[string]Response{},
			Description: "whoops forgot to add desc",
		}
		if paths[r.Path] == nil {
			paths[r.Path] = make(map[string]interface{})
		}
		paths[r.Path][strings.ToLower(r.Protocol)] = pathMethod
	}

	info := Info{
		Version: appVersion,
		Title:   appName,
	}
	openapi := OpenAPI{
		SpecVersion: "3.1.0",
		Info:        info,
		Servers:     []Server{},
		//FIX ME: SERVERS SHOULD BE SPECIFIED IN THE APP CONFIGURATION
		//FIX ME: THE DEFAULT SERVER SHOULD BE THE NETWORK IP: PORT
		Paths: paths,
	}
	openapiJSON, err := json.Marshal(openapi)
	if err != nil {
		return "", err
	}
	return string(openapiJSON), nil
}
