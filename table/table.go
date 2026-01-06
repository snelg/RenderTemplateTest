package table

import (
	"embed"
	"html/template"
)

//go:embed templates
var tableFiles embed.FS

type TableConfig struct {
	Headers []Header
	Class   string
}

type Header struct {
	Label string
	// Other options will go here, e.g. `Sortable bool`
}

func AddTableTemplate(sourceTemplate *template.Template) (*template.Template, error) {
	_, err := sourceTemplate.ParseFS(tableFiles, "templates/table.tmpl")

	return sourceTemplate, err
}
