package utils

import (
	"bytes"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//BuildURL - Builds a URL in https://www.example.com/ format
func BuildURL(domain string, path string) string {
	log.Println("base domain", getBaseDomain(domain))
	if strings.HasPrefix(path, "http") {
		return path
	}

	// return ancor link
	if !strings.HasPrefix(path, "http") && strings.HasPrefix(path, "#") {
		return domain + path
	}

	domain = getBaseDomain(domain)
	//verify if path is not an url and if, for some reason, doesn't starts with a /
	if !strings.HasPrefix(path, "http") && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return domain + path
}

func getBaseDomain(s string) string {
	url, err := url.Parse(FormatURL(s))
	if err != nil {
		log.Println(err)
	}
	domain := url.Hostname()
	return FormatURL(domain)
}

//GetLinks - returns an array with all links found in the page
func GetLinks(body *html.Node, url string) []string {
	var urls []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			//TODO just works if href appears just before a tag
			for _, link := range n.Attr {
				if link.Key == "href" {
					//TODO change to get domain
					url := BuildURL(url, link.Val)
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

//GetHTMLElement - returns an array with all html element found in the page
func GetHTMLElement(body *html.Node, htmlElement string) []string {
	var element *html.Node
	var stringElements []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == htmlElement {
			element = n
			stringElements = append(stringElements, formatHTML(element))
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(body)
	return stringElements
}

//FormatURL - adds http to the beginning of an URL it does not contain
func FormatURL(url string) string {
	if strings.Contains(url, "http://") || strings.Contains(url, "https://") {
		return url
	}
	return "http://" + url
}

func formatHTML(element *html.Node) string {
	var buffer bytes.Buffer
	w := io.Writer(&buffer)
	html.Render(w, element)
	return buffer.String()
}

//ReadBody - return an array of bytes of the resp.Body
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

//ParseBody - returns a *html.Node of an array of bytes of the resp.Body
func ParseBody(body []byte) *html.Node {
	bodyParsed, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Println(err)
	}
	return bodyParsed
}
