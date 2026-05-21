package docs

import "embed"

//go:embed docs.html openapi.yaml
var FS embed.FS
