package main

import (
	"github.com/abyssparanoia/rapid-go/cmd/helper/handler"
	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
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
	Logger        *log.Logger
	HelperHandler handler.HelperHandler
}

// Inject ... indect dependency
func (d *Dependency) Inject(e *Environment) {

	var lCli log.Writer

	if e.ENV == "LOCAL" {
		lCli = log.NewWriterStdout()
	} else {
		lCli = log.NewWriterStackdriver(e.ProjectID)
	}

	d.Logger = log.NewLogger(lCli, log.NewSeverity(e.MinLogSeverity), "script")
	d.HelperHandler = handler.NewHelperHandler()

}
