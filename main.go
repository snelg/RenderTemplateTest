package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
)

type User struct {
	Name  string
	Email string
}
type Order struct {
	ID    string
	Total string
}

type TableData[T any] struct {
	Headers    []string
	Class      string
	TemplateID string
	Rows       []T
}

var tpl *template.Template // shared template set

// render executes a named template from THIS SAME template set and returns safe HTML.
func render(name string, data any) (template.HTML, error) {
	var b bytes.Buffer
	if err := tpl.ExecuteTemplate(&b, name, data); err != nil {
		return "", err
	}
	// Safe because html/template produced the markup (already escaped).
	return template.HTML(b.String()), nil
}

func mustParseTemplates() *template.Template {
	funcs := template.FuncMap{"render": render}
	// Parse component + layout first, then the page so its defines override blocks.
	tpl = template.Must(template.New("").Funcs(funcs).ParseFiles(
		"templates/components/table.tmpl",
		"templates/layouts/base.tmpl",
		"templates/pages/report.tmpl",
	))
	return tpl
}

func main() {
	mustParseTemplates()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			UsersTable  TableData[User]
			OrdersTable TableData[Order]
		}{
			UsersTable: TableData[User]{
				Headers:    []string{"Name", "Email"},
				Class:      "striped",
				TemplateID: "rows.users.slot",
				Rows:       []User{{"Ada", "ada@example.com"}, {"Linus", "linus@example.com"}},
			},
			OrdersTable: TableData[Order]{
				Headers:    []string{"Order #", "Total"},
				TemplateID: "rows.orders.slot",
				Rows:       []Order{{"A123", "$25.99"}, {"B456", "$79.99"}},
			},
		}
		// Execute the PAGE entrypoint; the page calls the layout internally.
		if err := tpl.ExecuteTemplate(w, "report", data); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
