FROM golang:1.11

COPY . /work
WORKDIR /work
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o url-analyzer

FROM scratch
COPY --from=0 /work/url-analyzer .
COPY --from=0 /work/templates/ /templates

EXPOSE 8080
ENTRYPOINT ["./url-analyzer"]
