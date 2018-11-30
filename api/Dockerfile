FROM golang:1.11-alpine3.8

COPY . /go/src/github.com/abyssparanoia/gke-beego/api
WORKDIR /go/src/github.com/abyssparanoia/gke-beego/api/src

RUN apk update \
    && apk add --no-cache git \
    && go get -u github.com/codegangsta/gin \
    && go get -u github.com/golang/dep/cmd/dep \
    && dep ensure

CMD gin -i run main.go dependency.go routing.go
