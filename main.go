package main

import (
	"github.com/KarineValenca/URL-analyzer/info"
	"html/template"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	var webpage info.WebPage
	webpage.URL = r.FormValue("url")
	if r.Method == http.MethodPost {
		webpage = info.BuildWebPageInfo(webpage)
	}

	tpl.ExecuteTemplate(w, "index.gohtml", webpage)
}
