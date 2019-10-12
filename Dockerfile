FROM golang:1.12-alpine

WORKDIR /go/src/github.com/abyssparanoia/rapid-go/
COPY . .

ENV GO111MODULE=off

RUN apk --no-cache --update upgrade \
    && apk add --no-cache git alpine-sdk \
    && go get -u github.com/pilu/fresh

ENV GO111MODULE=on

ENV PORT 8080
EXPOSE 8080

CMD fresh