FROM golang:1.19-alpine AS builder

WORKDIR /go/src/github.com/playground-live/moala-meet-and-greet-back/

ENV CGO_ENABLED=0

COPY . .

RUN make build