package main

import (
	"bytes"
	"fmt"
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

var tpl *template.Template // shared template set

// params builds a map[string]any so you can pass "named" args from templates.
func params(v ...any) (map[string]any, error) {
	if len(v)%2 != 0 {
		return nil, fmt.Errorf("params: odd number of arguments")
	}
	m := make(map[string]any, len(v)/2)
	for i := 0; i < len(v); i += 2 {
		k, ok := v[i].(string)
		if !ok {
			return nil, fmt.Errorf("params: key %d is not a string", i)
		}
		m[k] = v[i+1]
	}
	return m, nil
}

// slice lets you write (slice "A" "B") inside templates.
func slice(v ...string) []string { return v }

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
	funcs := template.FuncMap{
		"params": params,
		"slice":  slice,
		"render": render,
	}
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
			Users  []User
			Orders []Order
		}{
			Users:  []User{{"Ada", "ada@example.com"}, {"Linus", "linus@example.com"}},
			Orders: []Order{{"A123", "$25.99"}, {"B456", "$79.99"}},
		}
		// Execute the PAGE entrypoint; the page calls the layout internally.
		if err := tpl.ExecuteTemplate(w, "report", data); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
