package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
)

var templates *template.Template

type info struct {
	Ticker   string
	Date     string
	PrevDate string
	NextDate string
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

	date := r.URL.Query().Get("date")
	if date == "" {
		date = "2016-Mar-18"
	}

	const shortForm = "2006-Jan-02"
	max, _ := time.Parse(shortForm, "2016-Mar-18")

	t, _ := time.Parse(shortForm, date)
	prevDate := t.AddDate(0, 0, -1)
	if prevDate.Weekday() == 0 {
		fmt.Println("0")
		prevDate = t.AddDate(0, 0, -3)
	}

	nextDate, _ := time.Parse(shortForm, date)
	if !nextDate.Equal(max) {
		nextDate = nextDate.AddDate(0, 0, 1)
		if nextDate.Weekday() == 7 {
			nextDate = nextDate.AddDate(0, 0, 3)
		}
	}

	var prev string
	if prevDate.Day() < 10 {
		prev = fmt.Sprintf("%d-%s-0%d", prevDate.Year(), prevDate.Month().String()[:3], prevDate.Day())
	} else {
		prev = fmt.Sprintf("%d-%s-%d", prevDate.Year(), prevDate.Month().String()[:3], prevDate.Day())
	}

	var next string
	if nextDate.Day() < 10 {
		next = fmt.Sprintf("%d-%s-0%d", nextDate.Year(), nextDate.Month().String()[:3], nextDate.Day())
	} else {
		next = fmt.Sprintf("%d-%s-%d", nextDate.Year(), nextDate.Month().String()[:3], nextDate.Day())
	}

	i := info{
		strings.ToUpper(r.URL.Query().Get("ticker")),
		date,
		prev,
		next,
	}

	templates.ExecuteTemplate(w, "main.tmpl", i)
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
