package utils

import (
	"strings"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"log"
	"bytes"
	"io"
)

func buildUrl(domain string, path string) string {
	if strings.HasSuffix(domain, "/") {
		domain = domain[:len(domain)-len("/")]
	}
	if strings.Contains(path, "http://") || strings.Contains(path, "https://"){
		return path
	} else {
		return domain+path
	}
}

func getLinks(body *html.Node, url string) []string{
	var urls []string
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
	f(body)
	return urls
}

func getHtmlElement(body *html.Node, htmlElement string) []string {
	var element *html.Node
	var stringElements []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == htmlElement {
			element = n
			log.Println(element)
			stringElements = append(stringElements, formatHtml(element))
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(body)
	return stringElements
}

func formatHtml(element *html.Node) string{
	var buffer bytes.Buffer
	w := io.Writer(&buffer)
	html.Render(w, element)
	return buffer.String()
}

func ReadBody(resp *http.Response) []byte {
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

