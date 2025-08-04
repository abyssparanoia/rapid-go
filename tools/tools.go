//go:build tools

package tool

import (
	_ "github.com/aarondl/sqlboiler/v4"                         //nolint
	_ "github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-mysql" //nolint
	_ "github.com/air-verse/air"                                //nolint
	_ "github.com/bufbuild/buf/cmd/buf"                         //nolint
	_ "github.com/cloudspannerecosystem/wrench"                 //nolint
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"     //nolint
	_ "github.com/kauche/splanter"                              //nolint
	_ "go.mercari.io/yo"                                        //nolint
	_ "go.uber.org/mock/mockgen"                                //nolint
	_ "golang.org/x/tools/cmd/goimports"                        //nolint
	_ "mvdan.cc/gofumpt"                                        //nolint
)
