FROM golang:1.13-alpine

WORKDIR /go/src/github.com/abyssparanoia/rapid-go/
COPY . .

ENV GO111MODULE=off

RUN apk --no-cache --update upgrade \
    && apk add --no-cache git alpine-sdk \
    && go get github.com/tockins/realize

ENV GO111MODULE=on