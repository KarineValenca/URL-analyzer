package info

import (
	"github.com/KarineValenca/URL-analyzer/utils"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"regexp"
	"strings"
)

//WebPage - a struct to hold webpage info
type WebPage struct {
	URL                      string
	HTMLVersion              string
	PageTitle                string
	Headings                 heading
	CounterInternalLinks     int
	CounterExternalLinks     int
	CounterInaccessibleLinks int
	ContainsLoginForm        bool
	Error                    string
}

type heading struct {
	Counterh1 int
	Counterh2 int
	Counterh3 int
	Counterh4 int
	Counterh5 int
}

//BuildWebPageInfo - Build the information about the webpage
func BuildWebPageInfo(webpage WebPage) WebPage {
	resp, err := http.Get(utils.FormatURL(webpage.URL))
	if err != nil {
		log.Println(err)
		webpage.Error = "Invalid URL: try again"
	}

	if resp != nil {
		body := utils.ReadBody(resp)
		bodyParsed := utils.ParseBody(body)
		webpage.HTMLVersion = checkHTMLVersion(body)
		webpage.PageTitle = getPageTitle(bodyParsed)
		webpage.Headings.Counterh1 = len(utils.GetHTMLElement(bodyParsed, "h1"))
		webpage.Headings.Counterh2 = len(utils.GetHTMLElement(bodyParsed, "h2"))
		webpage.Headings.Counterh3 = len(utils.GetHTMLElement(bodyParsed, "h3"))
		webpage.Headings.Counterh4 = len(utils.GetHTMLElement(bodyParsed, "h4"))
		webpage.Headings.Counterh5 = len(utils.GetHTMLElement(bodyParsed, "h5"))
		webpage.CounterInternalLinks, webpage.CounterExternalLinks = countLinks(utils.GetHTMLElement(bodyParsed, "a"))
		webpage.CounterInaccessibleLinks = countInaccessibleLinks(utils.GetLinks(bodyParsed, webpage.URL))
		webpage.ContainsLoginForm = checkLoginFormPresence(bodyParsed)
	}
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
	titles := utils.GetHTMLElement(body, "title")
	if len(titles) > 0 {
		return titles[0]
	}

	return "Page has no title"

}

func countLinks(s []string) (int, int) {
	externalLinks := 0
	internalLinks := 0
	for i := range s {
		if strings.Contains(s[i], "http") {
			externalLinks++
		} else {
			internalLinks++
		}
	}
	return internalLinks, externalLinks
}

func countInaccessibleLinks(urls []string) int {
	inaccessibleLinks := 0
	for i := range urls {
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

func checkLoginFormPresence(body *html.Node) bool {
	containsEmail := false
	containsPassword := false
	inputs := utils.GetHTMLElement(body, "input")
	for i := range inputs {
		if strings.Contains(inputs[i], "email") || strings.Contains(inputs[i], "username") {
			containsEmail = true
		}
		if strings.Contains(inputs[i], "password") {
			containsPassword = true
		}
	}

	if containsEmail && containsPassword {
		return true
	}

	return false

}
