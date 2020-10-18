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
	webpage.HTMLVersion = checkHTMLVersion(body)
	webpage.PageTitle = getPageTitle(body)
	webpage.Headings.Counterh1 = len(getHtmlElement(body, "h1"))
	webpage.Headings.Counterh2 = len(getHtmlElement(body, "h2"))
	webpage.Headings.Counterh3 = len(getHtmlElement(body, "h3"))
	webpage.Headings.Counterh4 = len(getHtmlElement(body, "h4"))
	webpage.Headings.Counterh5 = len(getHtmlElement(body, "h5"))
	webpage.CounterInternalLinks, webpage.CounterExternalLinks = countLinks(getHtmlElement(body, "a"))
	webpage.CounterInaccessibleLinks = countInaccessibleLinks(getLinks(body, webpage.Url))
	webpage.ContainsLoginForm = checkLoginFormPresence(body)
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

func getPageTitle(body []byte) string {
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

func getLinks(body []byte, url string) []string{
	var urls []string
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Println(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			//TODO just works if href appears just before a tag
			for _, link := range n.Attr {
				if link.Key == "href" {
					//TODO change to get domain
					url := buildUrl(url, link.Val)
					urls = append(urls, url)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return urls
}

func buildUrl(domain string, url string) string {
	if strings.HasSuffix(domain, "/") {
		domain = domain[:len(domain)-len("/")]
	}
	if strings.Contains(url, "http://") || strings.Contains(url, "https://"){
		return url
	} else {
		return domain+url
	}
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

func checkLoginFormPresence(body []byte) bool{
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

//TODO add err
func getHtmlElement(body []byte, htmlElement string) []string {
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Println(err)
	}
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