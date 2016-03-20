package main

import (
	"html/template"
	"net/http"
	"strings"
)

var templates *template.Template

type info struct {
	Ticker string
}

func compileTemplates() {
	t, err := template.ParseFiles(
		"tmpl/header.tmpl",
		"tmpl/main.tmpl",
		"tmpl/index.tmpl")

	templates = template.Must(t, err)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.tmpl", nil)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "header.tmpl", nil)
	templates.ExecuteTemplate(w, "main.tmpl", &info{strings.ToUpper(r.URL.Query().Get("ticker"))})
}

func main() {
	compileTemplates()
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/main", mainHandler)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/imgs/", http.StripPrefix("/imgs/", http.FileServer(http.Dir("imgs"))))

	http.ListenAndServe(":8080", nil)
}
