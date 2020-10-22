[![Go Report Card](https://goreportcard.com/badge/github.com/KarineValenca/URL-analyzer)](https://goreportcard.com/report/github.com/KarineValenca/URL-analyzer)
[![Actions Status](https://github.com/KarineValenca/URL-analyzer/workflows/build/badge.svg)](https://github.com/KarineValenca/URL-analyzer/actions)
[![Actions Status](https://github.com/KarineValenca/URL-analyzer/workflows/test/badge.svg)](https://github.com/KarineValenca/URL-analyzer/actions)


# URL Analyzer
A web application to analyze web pages. The following information is provided:

- HTML Version: The HTML version of the web page (eg. `HTML 5 Doctype`, `HTML 4.01 Strict`, `XHTML 1.1`)
- Page Title: The `<title>` of the page.
- Counter h1: Counter the quantity of `<h1>` tags.
- Counter h2: Counter the quantity of `<h2>` tags.
- Counter h3: Counter the quantity of `<h3>` tags.
- Counter h4: Counter the quantity of `<h4>` tags.
- Counter h4: Counter the quantity of `<h5>` tags.
- Counter internal links: Counter the quantity of subpath links (eg: `/docs`).
- Counter external links: Counter the quantity of links with `http` or `https`.
- Counter inaccessible links: Counter of the quantity of links that returned a 400 or a 500 status, or that returned an error.
- Contains Login form: Checks if the page includes `<inputs>` with labels `email` or `username`, AND `password`. 

## Installation

You can run this app in two ways:

### Docker 
1. Clone this project

2. Build the docker image:

`docker build -t url-analyzer .`

3. Run the image:

`docker run -p 8080:8080 url-analyzer`

4. Access `localhost:8080`

### Manual installation
1. Clone this project and run:

`go run main.go`

2. Access `localhost:8080`

## Usage

Insert the URL and submit:

![Form](https://github.com/KarineValenca/URL-analyzer/blob/master/assets/image1.png
)

Wait some seconds until the app finishes the requests to the links, and you will see the result:

![Result](https://github.com/KarineValenca/URL-analyzer/blob/master/assets/image2.png)

> :warning: **NOTE**: 
> Depending on the amount of links on the webpage, the application will take more time to return the results, since a request is made to each link on the webpage.

## Tests

To execute the tests, run the command:

`go test ./...` 

To see the test coverage, run the command:

`go test -cover ./...`