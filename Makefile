SQLBOILER_SED_EXPRESSION := "s/{{GOPATH}}/$(subst /,\/,$(GOPATH))/g"


.PHONY: build
build:
	go build -o ./.bin/app-cli ./cmd/app

.PHONY: test
test:
	@go test ./internal/...

.PHONY: http.dev
http.dev:
	@go run github.com/air-verse/air -c .air.toml

.PHONY: generate.mock
generate.mock:
	@go generate ./...
	$(call format)

.PHONY: generate.buf
generate.buf:
	rm -rf ./internal/infrastructure/grpc/pb
	@go run github.com/bufbuild/buf/cmd/buf generate
	find ./schema/openapi/rapid/admin_api/v1 -type f ! -name 'api.swagger.json' -delete
	find ./schema/openapi/rapid/public_api/v1 -type f ! -name 'api.swagger.json' -delete
	$(call format)

.PHONY: generate.sqlboiler
generate.sqlboiler:
	@sed -e $(SQLBOILER_SED_EXPRESSION) ./db/main/sqlboiler.toml.tpl > ./db/main/sqlboiler.toml
	@go run github.com/volatiletech/sqlboiler/v4 --config=./db/main/sqlboiler.toml mysql
	@rm ./db/main/sqlboiler.toml
	$(call format)

.PHONY: generate.yo
generate.yo:
	@go run go.mercari.io/yo \
		$(SPANNER_PROJECT_ID) \
		$(SPANNER_INSTANCE_ID) \
		$(SPANNER_DATABASE_ID) \
		--out internal/infrastructure/spanner/internal/dbmodel \
		--package dbmodel \
		--template-path db/spanner/yo/templates \
		--suffix .go \
		--ignore-tables "SchemaMigrations"
	$(call format)

.PHONY: lint.go
lint.go:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run
	$(call format)

.PHONY: lint.go.fix
lint.go.fix:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run --fix
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
	.bin/app-cli schema-migration database sync-constants
	.bin/app-cli schema-migration database extract-schema

.PHONY: migrate.spanner.up
migrate.spanner.up:
	@go run github.com/cloudspannerecosystem/wrench migrate up --directory ./db/spanner
	@go run github.com/cloudspannerecosystem/wrench load --directory ./db/spanner
	@go run github.com/kauche/splanter \
		--project $(SPANNER_PROJECT_ID) \
		--instance $(SPANNER_INSTANCE_ID) \
		--database $(SPANNER_DATABASE_ID) \
		--directory ./db/spanner/masterdata

.PHONY: format
format:
	$(call format)

.PHONY: init.local.cognito
init.local.cognito:
	bash ./localstack/cognito/init.sh

define format
	@go fmt ./... 
	@go run github.com/bufbuild/buf/cmd/buf format -w
	@go run golang.org/x/tools/cmd/goimports -w ./ 
	@go run mvdan.cc/gofumpt -l -w .
	@go mod tidy
endef