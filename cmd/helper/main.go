package main

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"

	"github.com/abyssparanoia/rapid-go/cmd/helper/ctxhelper"
)

func main() {

	e := &Environment{}
	e.Get()

	// Dependency
	d := Dependency{}
	d.Inject(e)
	ctx := context.Background()
	ctx = ctxzap.ToContext(ctx, d.Logger)

	setDeps(d)
	ctxhelper.SetContext(ctx)

	execute()
}
