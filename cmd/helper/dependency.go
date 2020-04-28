package main

import (
	"github.com/abyssparanoia/rapid-go/cmd/helper/handler"
	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
	"go.uber.org/zap"
)

// Deps ...
var Deps Dependency

func setDeps(deps Dependency) {
	Deps = deps
}

func getDeps() Dependency {
	return Deps
}

// Dependency ... dependency
type Dependency struct {
	Logger        *zap.Logger
	HelperHandler handler.HelperHandler
}

// Inject ... indect dependency
func (d *Dependency) Inject(e *Environment) {

	d.Logger, _ = log.New("LOCAL")
	d.HelperHandler = handler.NewHelperHandler()

}
