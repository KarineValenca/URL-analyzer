package main

import (
	"strings"
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