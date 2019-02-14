FROM golang:1.11-alpine3.8

COPY . /go/src/github.com/abyssparanoia/rapid-go/api
WORKDIR /go/src/github.com/abyssparanoia/rapid-go/api/src

RUN apk --no-cache --update upgrade \
    && apk add --no-cache git \
    && go get -u github.com/codegangsta/gin \
    && go get -u github.com/golang/dep/cmd/dep \
    && dep ensure

CMD gin -i run main.go dependency.go routing.go
