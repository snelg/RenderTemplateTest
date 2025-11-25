package table

import (
	"embed"
	"html/template"
)

//go:embed templates
var tableFiles embed.FS

type TableData[T any] struct {
	Headers []Header
	Class   string
	Rows    []T
}

type Header struct {
	Label string
	// Other options will go here, e.g. `Sortable bool`
}

func AddTableTemplate(sourceTemplate *template.Template) (*template.Template, error) {
	_, err := sourceTemplate.ParseFS(tableFiles, "templates/table.tmpl")

	return sourceTemplate, err
}
