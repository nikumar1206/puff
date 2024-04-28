package openapi

import "fmt"

type OpenAPI struct {
	title string // title of the application
}

var OPENAPI_UI string = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <title>Elements in HTML</title>
  
    <script src="https://unpkg.com/@stoplight/elements/web-components.min.js"></script>
    <link rel="stylesheet" href="https://unpkg.com/@stoplight/elements/styles.min.css">
  </head>
  <body style="height: 100vh;">

    <elements-api
      apiDescriptionUrl="%.yaml"
      router="hash"
      tryItCredentialsPolicy="same-origin"
    />

  </body>
</html>`

func GenerateOpenAPIUI(path, title string) string {
	return fmt.Sprintf(OPENAPI_UI, path, title)
}
