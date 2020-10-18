package utils

import (
	"testing"
	"strings"
	"golang.org/x/net/html"
	"github.com/stretchr/testify/assert"
	"net/http"
)

func TestBuildUrlDomainWithSlashUrlPath(t *testing.T) {
	// test for domain with slash and url as a path
	domain := "https://golang.org/"
	url := "/doc/"

	result := BuildUrl(domain, url)

	assert.Equal(t, result, "https://golang.org/doc/")
	
}

func TestBuildUrlDomainWithoutSlashUrlPath(t *testing.T) {
	// test for domain without slash and url as a path
	domain := "https://golang.org"
	url := "/doc/"

	result := BuildUrl(domain, url)
	assert.Equal(t, result, "https://golang.org/doc/")
}

func TestBuildUrlUrlHttpsLink(t *testing.T) {
	// test for url as a https link
	domain := "https://golang.org/"
	url := "https://golang.org/"

	result := BuildUrl(domain, url)
	assert.Equal(t, result, "https://golang.org/")
}

func TestBuildUrlUrlHttpLink(t *testing.T) {
	// test for url as a http link
	domain := "http://golang.org/"
	url := "http://golang.org/"

	result := BuildUrl(domain, url)
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
	bodyParsed, err := html.Parse(strings.NewReader(body))
	if err != nil {
		t.Log(err)
	}
	url := "https://golang.org/"
	result := GetLinks(bodyParsed, url)
	
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
	bodyParsed, err := html.Parse(strings.NewReader(body))
	if err != nil {
		t.Log(err)
	}
	url := "https://golang.org/"
	result := GetLinks(bodyParsed, url)
	
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
	bodyParsed, err := html.Parse(strings.NewReader(body))
	if err != nil {
		t.Log(err)
	}
	url := "https://golang.org/"
	result := GetLinks(bodyParsed, url)
	
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
	bodyParsed, err := html.Parse(strings.NewReader(body))
	if err != nil {
		t.Log(err)
	}
	htmlElement := "button"

	result := GetHtmlElement(bodyParsed, htmlElement)
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
	bodyParsed, err := html.Parse(strings.NewReader(body))
	if err != nil {
		t.Log(err)
	}
	htmlElement := "a"

	result := GetHtmlElement(bodyParsed, htmlElement)
	assert.Equal(t, len(result), 0)
}

func TestReadBody(t *testing.T) {
	resp, err := http.Get("https://golang.org/")
	if err != nil {
		t.Log(err)
	}
	result := ReadBody(resp)

	assert.Contains(t, string(result), "<title>The Go Programming Language</title>")
}