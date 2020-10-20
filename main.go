package main

import (
	"github.com/KarineValenca/URL-analyzer/info"
	"github.com/KarineValenca/URL-analyzer/utils"
	"html/template"
	"log"
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
	if r.Method == http.MethodPost {
		webpage.URL = utils.FormatURL(r.FormValue("url"))
		resp, err := http.Get(webpage.URL)
		if err != nil {
			log.Println(err)
			webpage.Error = "Invalid URL: try again"
		}
		if resp != nil {
			webpage = info.BuildWebPageInfo(webpage, resp)
		}
	}

	tpl.ExecuteTemplate(w, "index.gohtml", webpage)
}
