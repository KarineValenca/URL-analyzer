package main

import (
	"net/http"
	"io/ioutil"
	"strings"
	"golang.org/x/net/html"
	"bytes"
	"io"
	"html/template"
	"log"
	"regexp"
)

type WebPage struct {
	Url string
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
	var webpage WebPage
	if r.Method == http.MethodPost {
		webpage.Url = r.FormValue("url")
		resp, err := http.Get(webpage.Url)
		if err != nil {
			log.Println(err)
		}
		webpage = buildWebPageInfo(webpage, resp)
	}
	
	tpl.ExecuteTemplate(w, "index.gohtml", webpage)
}

func buildWebPageInfo(webpage WebPage, resp *http.Response) WebPage {
	body := readBody(resp)
	bodyParsed := parseBody(body)
	webpage.HTMLVersion = checkHTMLVersion(body)
	webpage.PageTitle = getPageTitle(bodyParsed)
	webpage.Headings.Counterh1 = len(getHtmlElement(bodyParsed, "h1"))
	webpage.Headings.Counterh2 = len(getHtmlElement(bodyParsed, "h2"))
	webpage.Headings.Counterh3 = len(getHtmlElement(bodyParsed, "h3"))
	webpage.Headings.Counterh4 = len(getHtmlElement(bodyParsed, "h4"))
	webpage.Headings.Counterh5 = len(getHtmlElement(bodyParsed, "h5"))
	webpage.CounterInternalLinks, webpage.CounterExternalLinks = countLinks(getHtmlElement(bodyParsed, "a"))
	webpage.CounterInaccessibleLinks = countInaccessibleLinks(getLinks(bodyParsed, webpage.Url))
	webpage.ContainsLoginForm = checkLoginFormPresence(bodyParsed)
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

func getPageTitle(body *html.Node) string {
	titles := getHtmlElement(body, "title")
	if len(titles) > 0 {
		return titles[0]
	} else {
		return "Page has no title"
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

func countInaccessibleLinks(urls []string) int {
	inaccessibleLinks := 0
	for i, _ := range urls {
		resp, err := http.Get(urls[i])
		if err != nil {
			log.Println(err)
			inaccessibleLinks++
			continue
		}
		errRegex := regexp.MustCompile(`(4..|5..)`)
		if errRegex.Match([]byte(resp.Status)) {
			inaccessibleLinks++
		}
	}
	return inaccessibleLinks
}

func checkLoginFormPresence(body *html.Node) bool{
	containsEmail := false
	containsPassword := false
	inputs := getHtmlElement(body, "input")
	for i, _ := range inputs {
		if strings.Contains(inputs[i], "email") || strings.Contains(inputs[i], "username") {
			containsEmail = true
		}
		if strings.Contains(inputs[i], "password") {
			containsPassword = true
		}
	}

	if containsEmail && containsPassword {
		return true
	} else {
		return false
	}
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

func parseBody(body []byte) *html.Node {
	bodyParsed, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Println(err)
	}
	return bodyParsed
}

func formatHtml(element *html.Node) string{
	var buffer bytes.Buffer
	w := io.Writer(&buffer)
	html.Render(w, element)
	return buffer.String()
}