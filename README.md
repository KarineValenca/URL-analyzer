[![Go Report Card](https://goreportcard.com/badge/github.com/KarineValenca/URL-analyzer)](https://goreportcard.com/report/github.com/KarineValenca/URL-analyzer)
[![Actions Status](https://github.com/KarineValenca/URL-analyzer/workflows/build/badge.svg)](https://github.com/KarineValenca/URL-analyzer/actions)
[![Actions Status](https://github.com/KarineValenca/URL-analyzer/workflows/test/badge.svg)](https://github.com/KarineValenca/URL-analyzer/actions)


# URL Analyzer
A web application to analyze web pages.

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

## Tests

To execute tests, run the command:

`go test ./...` 

To see the test coverage, run the command:

`go test -cover ./...`