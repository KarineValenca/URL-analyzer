package utils

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"net/http"
	"strings"
	"testing"
)

func TestBuildURLDomainWithSlashUrlPath(t *testing.T) {
	// test for domain with slash and url as a path
	domain := "https://golang.org/"
	url := "/doc/"

	result := BuildURL(domain, url)

	assert.Equal(t, result, "http://golang.org/doc/")

}

func TestBuildURLDomainWithoutSlashUrlPath(t *testing.T) {
	// test for domain without slash and url as a path
	domain := "https://golang.org"
	url := "/doc/"

	result := BuildURL(domain, url)
	assert.Equal(t, result, "http://golang.org/doc/")
}

func TestBuildURLWithSubPath(t *testing.T) {
	domain := "https://golang.org/pkg/io/ioutil/"
	url := "/doc/"
	result := BuildURL(domain, url)
	assert.Equal(t, result, "http://golang.org/doc/")
}

func TestBuildURLWithSubPathAndAncorLink(t *testing.T) {
	domain := "https://golang.org/pkg/io/ioutil/"
	url := "#example_WriteFile"

	result := BuildURL(domain, url)
	assert.Equal(t, result, "https://golang.org/pkg/io/ioutil/#example_WriteFile")
}

func TestBuildURLHttpsLink(t *testing.T) {
	// test for url as a https link
	domain := "https://golang.org/"
	url := "https://golang.org/"

	result := BuildURL(domain, url)
	assert.Equal(t, result, "https://golang.org/")
}

func TestBuildURLHttpLink(t *testing.T) {
	// test for url as a http link
	domain := "http://golang.org/"
	url := "http://golang.org/"

	result := BuildURL(domain, url)
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
	assert.Contains(t, result, "http://golang.org/doc/copyright.html")
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
	assert.Contains(t, result, "http://golang.org/doc/copyright.html")
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

func TestGetHTMLElements(t *testing.T) {
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

	result := GetHTMLElement(bodyParsed, htmlElement)
	assert.Equal(t, len(result), 2)
	assert.Contains(t, result, `<button class="Button js-playgroundShareEl" title="Share this code">Share</button>`)
	assert.Contains(t, result, `<button class="Button js-playgroundShareEl" title="Like this code">Like</button>`)
}

func TestGetHTMLElementsNoElement(t *testing.T) {
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

	result := GetHTMLElement(bodyParsed, htmlElement)
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

func TestFormatURLFullURL(t *testing.T) {
	url := "https://golang.org"
	result := FormatURL(url)

	assert.Equal(t, result, "https://golang.org")
}

func TestFormatURLIncompleteURL(t *testing.T) {
	url := "golang.org"
	result := FormatURL(url)

	assert.Equal(t, result, "http://golang.org")
}
