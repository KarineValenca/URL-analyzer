package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	"golang.org/x/net/html"
	"bytes"
	"io"
)

func main()  {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://golang.org/")
	if err != nil {
		fmt.Println(err)
	}
	
	getHtmlElement(resp, "h2")
	fmt.Fprintf(w, "Hello world")
}


func getHtmlElement(resp *http.Response, htmlElement string) {
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	doc, _ := html.Parse(strings.NewReader(string(body)))
	var element *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == htmlElement {
			element = n
			formatHtml(element)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}

func formatHtml(element *html.Node) {
	var buffer bytes.Buffer
	w := io.Writer(&buffer)
	html.Render(w, element)
	fmt.Println(buffer.String())
}