.PHONY: build
build:
	go build -o ./.bin/app-cli ./cmd/app


.PHONY: test
test:
	@go test ./internal/...

.PHONY: http.dev
http.dev:
	@go run github.com/cosmtrek/air -c .air.toml

.PHONY: generate.mock
generate.mock:
	@go generate ./...
	$(call format)

.PHONY: generate.buf
generate.buf:
	@go run github.com/bufbuild/buf/cmd/buf generate
	$(call format)

.PHONY: generate.sqlboiler
generate.sqlboiler:
	@go run github.com/volatiletech/sqlboiler/v4 --config=./db/main/sqlboiler.toml mysql
	$(call format)

.PHONY: lint.go
lint.go:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run
	$(call format)

.PHONY: lint.proto
lint.proto:
	@go run github.com/bufbuild/buf/cmd/buf lint
	$(call format)

.PHONY: migrate.up
migrate.up:
	@go run github.com/pressly/goose/v3/cmd/goose --dir db/main/migrations mysql "$(DB_USER):$(DB_PASSWORD)@$(DB_HOST)/$(DB_DATABASE)?parseTime=true" up

.PHONY: format
format:
	$(call format)

define format
	@go fmt ./... 
	@go run github.com/bufbuild/buf/cmd/buf format
	@go run golang.org/x/tools/cmd/goimports -w ./ 
	@go mod tidy
endef