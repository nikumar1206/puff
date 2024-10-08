package puff

var openAPIHTML string = `<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="theme-color" content="#fff" />
    <meta name="description" content="SwaggerUIMultifold" />
    <link rel="icon" type="image/x-icon" href="https://fav.farm/💨" />
    <link
      rel="stylesheet"
      href="https://unpkg.com/swagger-editor@5.0.0-alpha.86/dist/swagger-editor.css"
    />
    <title>%s</title>
    <h1 id="connection-status"></h1>
  </head>
  <body style="margin: 0; padding: 0">
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
        url: "%s",
        dom_id: "#swagger-ui",
        presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
        plugins: [
          SwaggerEditor.plugins.EditorContentType,
          SwaggerEditor.plugins.EditorPreviewAsyncAPI,
          SwaggerEditor.plugins.EditorPreviewApiDesignSystems,
          SwaggerEditor.plugins.SwaggerUIAdapter,
          SwaggerUIBundle.plugins.DownloadUrl,
        ],
        layout: "StandaloneLayout",
      });
      document.body.onload = function () {
        const sc = document.getElementsByClassName("swagger-container");
        const tb = document.getElementsByClassName("topbar");
        sc[0].removeChild(tb[0]);
      };
    </script>
  </body>
</html>
`
