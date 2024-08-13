FROM golang:1.23-alpine AS builder

WORKDIR /go/src/github.com/abyssparanoia/rapid-go/

ENV CGO_ENABLED=0

COPY . .

RUN make build