package main

import (
	"context"

	"github.com/abyssparanoia/rapid-go/cmd/helper/ctxhelper"
	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
)

func main() {

	e := &Environment{}
	e.Get()

	// Dependency
	d := Dependency{}
	d.Inject(e)
	ctx := context.Background()
	ctx = log.SetLogger(ctx, d.Logger)

	setDeps(d)
	ctxhelper.SetContext(ctx)

	execute()
}
