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

.PHONY: migrate.create
migrate.create:
	make build
	.bin/app-cli schema-migration database create

.PHONY: migrate.up
migrate.up:
	make build
	.bin/app-cli schema-migration database up

.PHONY: format
format:
	$(call format)

.PHONY: init.local.cognito
init.local.cognito:
	bash ./localstack/cognito/init.sh

define format
	@go fmt ./... 
	@go run github.com/bufbuild/buf/cmd/buf format
	@go run golang.org/x/tools/cmd/goimports -w ./ 
	@go mod tidy
endef