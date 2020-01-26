# note: call scripts from /scripts

test:
	go test `go list ./... | grep -v pkg/dbmodels`