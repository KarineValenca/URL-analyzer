package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	"golang.org/x/net/html"
	"bytes"
	"io"
	"html/template"
)

type WebPage struct {
	HTMLVersion string
	PageTitle string
	Headings []string
	CounterInternalLinks int
	CounterExternalLinks int
	CounterInaccessibleLinks int
	ContainsLoginForm bool
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}
func main()  {
	
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}



func index(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://golang.org/")
	if err != nil {
		fmt.Println(err)
	}
	
	var webpage WebPage
	webpage.PageTitle = getHtmlElement(resp, "title")[0] //todo
 	fmt.Println(webpage.PageTitle)
	tpl.ExecuteTemplate(w, "index.gohtml", webpage)
}


func getHtmlElement(resp *http.Response, htmlElement string) []string {
	//converts http.Response.Body to *html.Node
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	doc, _ := html.Parse(strings.NewReader(string(body)))

	//get html element
	var element *html.Node
	var stringElements []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == htmlElement {
			element = n
			stringElements = append(stringElements, formatHtml(element))
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return stringElements
}


func formatHtml(element *html.Node) string{
	var buffer bytes.Buffer
	w := io.Writer(&buffer)
	html.Render(w, element)
	return buffer.String()
}