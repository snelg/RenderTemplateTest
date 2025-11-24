package table

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed templates
var tableFiles embed.FS

type TableData struct {
	Headers []Header
	Class   string
	Rows    any
}

type tableDataWithID struct {
	TableData
	TemplateKey string
}

type Header struct {
	Label string
	// Other options will go here, e.g. `Sortable bool`
}

func Renderer(bodyTemplate *template.Template) func(string, TableData) (template.HTML, error) {
	tableTemplate := template.New("table.tmpl") // template.ParseFS uses the file's base name as the main key. So if you change the actual file name (see below), you'll also need to change it here
	tableTemplate.Funcs(template.FuncMap{"renderRows": rowsRenderer(bodyTemplate)})
	template.Must(tableTemplate.ParseFS(tableFiles, "templates/table.tmpl")) // See above comment about template.ParseFS

	return func(templateKey string, data TableData) (template.HTML, error) {
		var b bytes.Buffer
		wrappedData := tableDataWithID{
			TableData:   data,
			TemplateKey: templateKey,
		}
		if err := tableTemplate.Execute(&b, wrappedData); err != nil {
			return "", err
		}
		// Safe because html/template produced the markup (already escaped).
		return template.HTML(b.String()), nil
	}
}

func rowsRenderer(bodyTemplate *template.Template) func(string, any) (template.HTML, error) {
	return func(name string, data any) (template.HTML, error) {
		var b bytes.Buffer
		if err := bodyTemplate.ExecuteTemplate(&b, name, data); err != nil {
			return "", err
		}
		// Safe because html/template produced the markup (already escaped).
		return template.HTML(b.String()), nil
	}
}
