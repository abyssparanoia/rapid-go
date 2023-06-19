//go:build tools

package tool

import (
	_ "github.com/bufbuild/buf/cmd/buf"                              //nolint
	_ "github.com/cosmtrek/air"                                      //nolint
	_ "github.com/golang/mock/mockgen"                               //nolint
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"          //nolint
	_ "github.com/volatiletech/sqlboiler/v4"                         //nolint
	_ "github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql" //nolint
	_ "golang.org/x/tools/cmd/goimports"                             //nolint
)
