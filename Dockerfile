FROM golang:1.11

COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o url-analyzer

EXPOSE 8080
ENTRYPOINT ["./url-analyzer"]
