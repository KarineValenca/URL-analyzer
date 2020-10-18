package main

import (
	"testing"
	"strings"
	"golang.org/x/net/html"
	"github.com/stretchr/testify/assert"
)

func TestBuildUrlDomainWithSlashUrlPath(t *testing.T) {
	// test for domain with slash and url as a path
	domain := "https://golang.org/"
	url := "/doc/"

	result := buildUrl(domain, url)

	assert.Equal(t, result, "https://golang.org/doc/")
	
}

func TestBuildUrlDomainWithoutSlashUrlPath(t *testing.T) {
	// test for domain without slash and url as a path
	domain := "https://golang.org"
	url := "/doc/"

	result := buildUrl(domain, url)
	assert.Equal(t, result, "https://golang.org/doc/")
}

func TestBuildUrlUrlHttpsLink(t *testing.T) {
	// test for url as a https link
	domain := "https://golang.org/"
	url := "https://golang.org/"

	result := buildUrl(domain, url)
	assert.Equal(t, result, "https://golang.org/")
}

func TestBuildUrlUrlHttpLink(t *testing.T) {
	// test for url as a http link
	domain := "http://golang.org/"
	url := "http://golang.org/"

	result := buildUrl(domain, url)
	assert.Equal(t, result, "http://golang.org/")
}

func TestGetLinks(t *testing.T) {
	body := `
	<body class="Site">
		<button class="Button js-playgroundShareEl" title="Share this code">Share</button>
		<a href="https://tour.golang.org/" title="Playground Go from your browser">Tour</a>
		<ul class="Footer-links">
			  <li class="Footer-link"><a href="/doc/copyright.html">Copyright</a></li>
		</ul>
	</body>
	`
	bodyParsed, _ := html.Parse(strings.NewReader(body))
	url := "https://golang.org/"
	result := getLinks(bodyParsed, url)
	
	assert.Equal(t, len(result), 2)
	assert.Contains(t, result, "https://tour.golang.org/")
	assert.Contains(t, result, "https://golang.org/doc/copyright.html")
}

func TestGetLinksWithClass(t *testing.T) {
	//test getLinks a tag with class
	body := `
	<body class="Site">
		<button class="Button js-playgroundShareEl" title="Share this code">Share</button>
		<a class="Button tour" href="https://tour.golang.org/" title="Playground Go from your browser">Tour</a>
		<ul class="Footer-links">
			  <li class="Footer-link"><a href="/doc/copyright.html">Copyright</a></li>
		</ul>
	</body>
	`
	bodyParsed, _ := html.Parse(strings.NewReader(body))
	url := "https://golang.org/"
	result := getLinks(bodyParsed, url)
	
	assert.Equal(t, len(result), 2)
	assert.Contains(t, result, "https://tour.golang.org/")
	assert.Contains(t, result, "https://golang.org/doc/copyright.html")
}

func TestGetLinksNoLink(t *testing.T) {
	//test getLinks a tag with class
	body := `
	<body class="Site">
		<button class="Button js-playgroundShareEl" title="Share this code">Share</button>
		<ul class="Footer-links">
			  <li class="Footer-link"Copyright</li>
		</ul>
	</body>
	`
	bodyParsed, _ := html.Parse(strings.NewReader(body))
	url := "https://golang.org/"
	result := getLinks(bodyParsed, url)
	
	assert.Equal(t, len(result), 0)
}

func TestGetHtmlElements(t *testing.T) {
	body := `
	<body class="Site">
		<button class="Button js-playgroundShareEl" title="Share this code">Share</button>
		<button class="Button js-playgroundShareEl" title="Like this code">Like</button>
		<ul class="Footer-links">
			  <li class="Footer-link"Copyright</li>
		</ul>
	</body>
	`
	bodyParsed, _ := html.Parse(strings.NewReader(body))
	htmlElement := "button"

	result := getHtmlElement(bodyParsed, htmlElement)
	assert.Equal(t, len(result), 2)
	assert.Contains(t, result, `<button class="Button js-playgroundShareEl" title="Share this code">Share</button>`)
	assert.Contains(t, result, `<button class="Button js-playgroundShareEl" title="Like this code">Like</button>`)
}

func TestGetHtmlElementsNoElement(t *testing.T) {
	body := `
	<body class="Site">
		<button class="Button js-playgroundShareEl" title="Share this code">Share</button>
		<button class="Button js-playgroundShareEl" title="Like this code">Like</button>
		<ul class="Footer-links">
			  <li class="Footer-link"Copyright</li>
		</ul>
	</body>
	`
	bodyParsed, _ := html.Parse(strings.NewReader(body))
	htmlElement := "a"

	result := getHtmlElement(bodyParsed, htmlElement)
	assert.Equal(t, len(result), 0)
}