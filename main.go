package main

import (
	"html/template"
	"log"
	"net/http"

	"TableRenderTest/table"
)

type User struct {
	Name  string
	Email string
}
type Order struct {
	ID    string
	Total string
}

type PageData struct {
	UsersTable  table.TableData[User]
	OrdersTable table.TableData[Order]
}

func mustParseTemplates() *template.Template {
	tpl := template.New("")
	// Parse component + layout first, then the page so its defines override blocks.
	template.Must(tpl.ParseFiles(
		"templates/layouts/base.tmpl",
		"templates/pages/report.tmpl",
	))
	template.Must(table.AddTableTemplate(tpl))
	return tpl
}

func main() {
	tpl := mustParseTemplates()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			UsersTable: table.TableData[User]{
				Headers: []table.Header{{Label: "Name"}, {Label: "Email"}},
				Class:   "striped",
				Rows:    []User{{"Ada", "ada@example.com"}, {"Linus", "linus@example.com"}},
			},
			OrdersTable: table.TableData[Order]{
				Headers: []table.Header{{Label: "Order #"}, {Label: "Total"}},
				Rows:    []Order{{"A123", "$25.99"}, {"B456", "$79.99"}},
			},
		}
		// Execute the PAGE entrypoint; the page calls the layout internally.
		if err := tpl.ExecuteTemplate(w, "report.tmpl", data); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
