package main

import (
	"html/template"
	"net/http"
)

var templates *template.Template

func compileTemplates() {
	t, err := template.ParseFiles(
		"tmpl/header.tmpl",
		"tmpl/footer.tmpl",
		"tmpl/index.tmpl")

	templates = template.Must(t, err)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "header.tmpl", nil)
	templates.ExecuteTemplate(w, "index.tmpl", nil)
	templates.ExecuteTemplate(w, "footer.tmpl", nil)
}

func main() {
	compileTemplates()
	http.HandleFunc("/", rootHandler)
	http.ListenAndServe(":8080", nil)
}
