FROM golang:1.10

WORKDIR /go/src/github.com/pendolf/ariadna
COPY . .

RUN go install -v ./...

EXPOSE 8080
CMD ["ariadna", "http"]
