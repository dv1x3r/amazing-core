package api

import (
	"embed"
	"html/template"
)

//go:embed templates/*.gohtml
var templatesFS embed.FS

var tmpl = template.Must(template.ParseFS(templatesFS, "templates/*.gohtml"))
