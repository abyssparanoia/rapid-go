FROM golang:1.19-alpine AS builder

WORKDIR /go/src/github.com/abyssparanoia/rapid-go/

ENV CGO_ENABLED=0

COPY . .

RUN go install -v -tags netgo -ldflags '-extldflags "-static"' ./cmd/app
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine AS server

RUN apk add ca-certificates
COPY --from=builder /go/bin/app /bin/app
COPY --from=builder /go/bin/goose /bin/goose

WORKDIR /go/src/github.com/abyssparanoia/rapid-go/

ENV PORT 8080

EXPOSE 8080