FROM golang:1.13.11-alpine3.10 AS builder

ARG SERVICE_NAME=default

WORKDIR /go/src/github.com/abyssparanoia/rapid-go/

ENV CGO_ENABLED=0

COPY . .

RUN go install -v -tags netgo -ldflags '-extldflags "-static"' ./cmd/rapid/


FROM alpine AS server

RUN apk add ca-certificates
COPY --from=builder /go/bin/rapid /bin/rapid

WORKDIR /go/src/github.com/abyssparanoia/rapid-go/

ENV PORT 8080
EXPOSE 8080

ENTRYPOINT ["rapid","default-http","run"]