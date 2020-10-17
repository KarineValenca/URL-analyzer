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
	body := readBody(resp)
	webpage.HTMLVersion = checkHTMLVersion(body)
	webpage.PageTitle = getHtmlElement(body, "title")[0] //todo clear readness
	webpage.Headings.Counterh1 = len(getHtmlElement(body, "h1"))
	webpage.Headings.Counterh2 = len(getHtmlElement(body, "h2"))
	webpage.Headings.Counterh3 = len(getHtmlElement(body, "h3"))
	webpage.Headings.Counterh4 = len(getHtmlElement(body, "h4"))
	webpage.Headings.Counterh5 = len(getHtmlElement(body, "h5"))
	webpage.CounterInternalLinks, webpage.CounterExternalLinks = countLinks(getHtmlElement(body, "a"))
	return webpage
}

func checkHTMLVersion(body []byte) string {
	if strings.Contains(string(body), "<!DOCTYPE html>") {
		return "HTML5 doctype"
	} else if strings.Contains(string(body), "DTD HTML 4.01") {
		return "HTML 4.01 doctype"
	} else {
		return "Couldn't find HTML version"
	}
}

func countLinks(s []string) (int, int) {
	externalLinks := 0
	internalLinks := 0
	for i, _ := range s {
		if strings.Contains(s[i], "http") {
			externalLinks++
		} else{
			internalLinks++
		}
	}
	return internalLinks, externalLinks
}

func readBody(resp *http.Response) []byte {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) 
	if err != nil {
		log.Fatal(err)
	}

	//restore the io.readcloser to its original state
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body
}

//todo add err
func getHtmlElement(body []byte, htmlElement string) []string {
	doc, _ := html.Parse(strings.NewReader(string(body)))
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