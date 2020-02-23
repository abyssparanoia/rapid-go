# note: call scripts from /scripts

format:
	$(call format)

test:
	go test `go list ./... | grep -v pkg/dbmodels`

define format
	go fmt ./... && goimports -w ./ && go mod tidy

endef