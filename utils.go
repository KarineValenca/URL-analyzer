package main

import (
	"strings"
	"golang.org/x/net/html"
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