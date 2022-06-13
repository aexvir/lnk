package svc

import (
	"net/http"

	"github.com/aexvir/lnk/proto"
)

func OpenapiSchemaHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(proto.Schema)
}

var docspage = `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Lnk API</title>
    <script src="https://unpkg.com/@stoplight/elements/web-components.min.js"></script>
    <link rel="stylesheet" href="https://unpkg.com/@stoplight/elements/styles.min.css">
	<style>
	div[data-overlay-container="true"] {height: 100vh !important;}
	</style>
  </head>
  <body>
    <elements-api apiDescriptionUrl="/api/schema.yaml" router="hash" hideTryIt="true" />
  </body>
</html>
`

func OpenapiDocsHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(docspage))
}
