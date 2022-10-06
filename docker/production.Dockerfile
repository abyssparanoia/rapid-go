FROM golang:1.19-alpine AS builder

WORKDIR /go/src/github.com/playground-live/moala-meet-and-greet-back/

ENV CGO_ENABLED=0

COPY . .

RUN go install -v -tags netgo -ldflags '-extldflags "-static"' ./cmd/app


FROM alpine AS server

RUN apk add ca-certificates
COPY --from=builder /go/bin/app /bin/app

WORKDIR /go/src/github.com/playground-live/moala-meet-and-greet-back/

EXPOSE 80
