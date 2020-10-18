package main

import (
	"github.com/KarineValenca/URL-analyzer/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestBuildWebPageInfo(t *testing.T) {
	var webpage webPage
	webpage.URL = "https://golang.org/"
	resp, err := http.Get(webpage.URL)
	if err != nil {
		t.Log(err)
	}
	result := buildWebPageInfo(webpage, resp)
	assert.Equal(t, result.URL, "https://golang.org/")
	assert.Equal(t, result.HTMLVersion, "HTML5 doctype")
	assert.Equal(t, result.PageTitle, "<title>The Go Programming Language</title>")
	assert.Equal(t, result.Headings.Counterh1, 1)
	assert.Equal(t, result.Headings.Counterh2, 3)
	assert.Equal(t, result.Headings.Counterh3, 0)
	assert.Equal(t, result.Headings.Counterh4, 0)
	assert.Equal(t, result.Headings.Counterh5, 0)
	assert.Equal(t, result.CounterInternalLinks, 10)
	assert.Equal(t, result.CounterExternalLinks, 7)
	assert.Equal(t, result.CounterInaccessibleLinks, 0)
	assert.Equal(t, result.ContainsLoginForm, false)
}

func TestHtmlVersionHtml5(t *testing.T) {
	body := []byte(`<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>Index</title>
		</head>
		<body>
		<h1>Hello World!</h1>
		</body>`)

	result := checkHTMLVersion(body)
	assert.Equal(t, result, "HTML5 doctype")
}

func TestHtmlVersionHtml4(t *testing.T) {
	body := []byte(`<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>Index</title>
		</head>
		<body>
		<h1>Hello World!</h1>
		</body>`)

	result := checkHTMLVersion(body)
	assert.Equal(t, result, "HTML 4.01 doctype")
}

func TestHtmlVersionCantFindVersion(t *testing.T) {
	body := []byte(`
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>Index</title>
		</head>
		<body>
		<h1>Hello World!</h1>
		</body>`)

	result := checkHTMLVersion(body)
	assert.Equal(t, result, "Couldn't find HTML version")
}

func TestGetPageTitle(t *testing.T) {
	body := []byte(`<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>Index</title>
		</head>
		<body>
		<h1>Hello World!</h1>
		</body>`)
	bodyParsed := utils.ParseBody(body)
	result := getPageTitle(bodyParsed)
	assert.Equal(t, result, "<title>Index</title>")
}

func TestGetPageTitleNoTitle(t *testing.T) {
	body := []byte(`<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
		</head>
		<body>
		<h1>Hello World!</h1>
		</body>`)
	bodyParsed := utils.ParseBody(body)
	result := getPageTitle(bodyParsed)
	assert.Equal(t, result, "Page has no title")
}

func TestCountLinks(t *testing.T) {
	urls := []string{
		"https://golang.org",
		"http://golang.org",
		"/doc",
	}

	internalLink, externalLink := countLinks(urls)

	assert.Equal(t, internalLink, 1)
	assert.Equal(t, externalLink, 2)
}

func TestInaccessibleLinks(t *testing.T) {
	urls := []string{
		"https://golang.org",
		"http://invalidlink.org",
		"https://stackoverflow.com/questions/3009061",
	}

	result := countInaccessibleLinks(urls)

	assert.Equal(t, result, 2)
}

func TestCheckLoginFormPresence(t *testing.T) {
	body := []byte(`<form action="action_page.php" method="post">
		<div class="container">
		<label for="uname"><b>Username</b></label>
		<input type="text" placeholder="Enter Username" name="username" required>
	
		<label for="psw"><b>Password</b></label>
		<input type="password" placeholder="Enter Password" name="password" required>
	
		<button type="submit">Login</button>
		</div>
		<button type="button" class="cancelbtn">Cancel</button>
		</div>
	</form>`)
	bodyParsed := utils.ParseBody(body)
	result := checkLoginFormPresence(bodyParsed)

	assert.Equal(t, result, true)
}

func TestCheckLoginFormPresenceNoLoginForm(t *testing.T) {
	body := []byte(`<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
		</head>
		<body>
		<h1>Hello World!</h1>
		</body>`)
	bodyParsed := utils.ParseBody(body)
	result := checkLoginFormPresence(bodyParsed)

	assert.Equal(t, result, false)
}
