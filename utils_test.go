package main

import (
	"testing"
)

func TestBuildUrlDomainWithSlashUrlPath(t *testing.T) {
	// test for domain with slash and url as a path
	domain := "https://golang.org/"
	url := "/doc/"

	result := buildUrl(domain, url)
	if result != "https://golang.org/doc/" {
		t.Errorf("buildUrl(%v,%v) failed, expected %v, got %v", domain, url, "https://golang.org/doc/", result)
	}else {
		t.Logf("buildUrl(%v,%v) success, expected %v, got %v", domain, url, "https://golang.org/doc/", result)
	}
}

func TestBuildUrlDomainWithoutSlashUrlPath(t *testing.T) {
	// test for domain without slash and url as a path
	domain := "https://golang.org"
	url := "/doc/"

	result := buildUrl(domain, url)
	if result != "https://golang.org/doc/" {
		t.Errorf("buildUrl(%v,%v) failed, expected %v, got %v", domain, url, "https://golang.org/doc/", result)
	}else {
		t.Logf("buildUrl(%v,%v) success, expected %v, got %v", domain, url, "https://golang.org/doc/", result)
	}
}

func TestBuildUrlUrlHttpsLink(t *testing.T) {
	// test for url as a https link
	domain := "https://golang.org/"
	url := "https://golang.org/"

	result := buildUrl(domain, url)
	if result != "https://golang.org/" {
		t.Errorf("buildUrl(%v,%v) failed, expected %v, got %v", domain, url, "https://golang.org/", result)
	}else {
		t.Logf("buildUrl(%v,%v) success, expected %v, got %v", domain, url, "https://golang.org/", result)
	}
}

func TestBuildUrlUrlHttpLink(t *testing.T) {
	// test for url as a http link
	domain := "http://golang.org/"
	url := "http://golang.org/"

	result := buildUrl(domain, url)
	if result != "http://golang.org/" {
		t.Errorf("buildUrl(%v,%v) failed, expected %v, got %v", domain, url, "http://golang.org/", result)
	}else {
		t.Logf("buildUrl(%v,%v) success, expected %v, got %v", domain, url, "http://golang.org/", result)
	}
}

