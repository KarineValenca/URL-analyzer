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
	"log"
)

type WebPage struct {
	HTMLVersion string
	PageTitle string
	Headings Heading
	CounterInternalLinks int
	CounterExternalLinks int
	CounterInaccessibleLinks int
	ContainsLoginForm bool
}

type Heading struct {
	Counterh1 int
	Counterh2 int
	Counterh3 int
	Counterh4 int
	Counterh5 int
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
	webpage = buildWebPageInfo(webpage, resp)
	tpl.ExecuteTemplate(w, "index.gohtml", webpage)
}

func buildWebPageInfo(webpage WebPage, resp *http.Response) WebPage {
	webpage.PageTitle = getHtmlElement(resp, "title")[0] //todo clear readness
	webpage.Headings.Counterh1 = len(getHtmlElement(resp, "h1"))
	webpage.Headings.Counterh2 = len(getHtmlElement(resp, "h2"))
	webpage.Headings.Counterh3 = len(getHtmlElement(resp, "h3"))
	webpage.Headings.Counterh4 = len(getHtmlElement(resp, "h4"))
	webpage.Headings.Counterh5 = len(getHtmlElement(resp, "h5"))
	return webpage
}

//todo add err
func getHtmlElement(resp *http.Response, htmlElement string) []string {
	//converts http.Response.Body to *html.Node
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) 
	if err != nil {
		fmt.Println("Erro do body")
		log.Fatal(err)
	}

	doc, _ := html.Parse(strings.NewReader(string(body)))
	
	//restore the io.readcloser to its original state
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

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